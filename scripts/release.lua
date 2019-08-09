-- sudo luarocks install argparse luasec inspect

-- # TODO: always update asset and body with checksum together

local ltn12 = require "ltn12"
local json = require "cjson"
local argparse = require "argparse"
local sha2 = require "sha2"
local http = require "socket.http"
local https = require "ssl.https"
local inspect = require "inspect"

local parser = argparse()
parser:name "git-team-releaser"
parser:description "Interactively release git-team via github api v3"
parser:epilog "https://github.com/hekmekk/git-team"
parser:option "--github-api-token"
parser:option "--git-team-version"
parser:option "--git-team-deb-path"

function find_release_in_tags(releases_uri, version)
  local respbody = {}

  local result, respcode, respheaders, respstatus = https.request {
      method = 'GET',
      url = releases_uri .. '/tags/' .. version,
      sink = ltn12.sink.table(respbody)
  }

  respbody = table.concat(respbody)

  if respbody and respbody ~= "" then
    return respcode, json.decode(respbody)
  end

  return respcode, nil
end

function upload_asset(github_api_token, upload_url, deb_file)
  local respbody = {} -- for the response body

  local result, respcode, respheaders, respstatus = https.request {
      method = 'POST',
      url = upload_url,
      source = ltn12.source.string(deb_file),
      headers = {
          ['Content-Type'] = 'application/vnd.debian.binary-package',
          ['Content-Length'] = tostring(#deb_file),
          ['Authorization'] = string.format('token %s', github_api_token)
      },
      sink = ltn12.sink.table(respbody)
  }
  respbody = table.concat(respbody)

  if respcode == 201 then
    print('[debug] successfully uploaded asset')
    return nil
  end

  if respbody and respbody ~= "" then
    return json.decode(respbody)
  end

  return nil
end

function update_release(github_api_token, releases_uri, release_id, request_body)
  local respbody = {} -- for the response body

  local result, respcode, respheaders, respstatus = https.request {
      method = 'PUT',
      url = releases_uri .. '/' .. release_id,
      source = ltn12.source.string(request_body),
      headers = {
          ['Content-Type'] = 'application/json',
          ['Content-Length'] = tostring(#request_body),
          ['Authorization'] = string.format('token %s', github_api_token)
      },
      sink = ltn12.sink.table(respbody)
  }
  respbody = table.concat(respbody)

  if respbody and respbody ~= "" then
    return respcode, json.decode(respbody)
  end

  return respcode, nil
end

function create_release(releases_uri, github_api_token, version)
  local reqbody = json.encode({ tag_name = version })
  local respbody = {} -- for the response body

  local result, respcode, respheaders, respstatus = https.request {
      method = 'POST',
      url = releases_uri,
      source = ltn12.source.string(reqbody),
      headers = {
          ['Content-Type'] = 'application/json',
          ['Content-Length'] = tostring(#reqbody),
          ['Authorization'] = string.format('token %s', github_api_token)
      },
      sink = ltn12.sink.table(respbody)
  }
  respbody = table.concat(respbody)

  if respbody and respbody ~= "" then
    return respcode, json.decode(respbody)
  end

  return respcode, nil
end

function is_checksum_in_body_already(release, checksum)
  if release and release['body'] == string.format('**sha256 checksum:** `%s`', checksum) then
    return true
  end
  return false
end

function is_asset_uploaded_already(assets, deb_file_name)
  for _, asset in pairs(assets) do
    if asset['name'] and asset['name'] == deb_file_name then
      return true
    end
    return false
  end
end

function read_file(file)
    local f = assert(io.open(file, "rb"))
    local content = f:read("*all")
    f:close()
    return content
end

function interactively_add_git_tag_and_push_to_remote(version)
  local x = os.execute(string.format('git tag -a %s', version))
  if 0 == x then
    print(string.format('[debug] latest commit has been tagged', version))
    os.execute('git push origin --tags')
    print(string.format('[debug] remote tags have been updated', version))
  end
end

-- main program

local args = parser:parse()

local github_api_token = args['github_api_token']
local git_team_version = args['git_team_version']
local git_team_deb_path = args['git_team_deb_path']
local releases_uri = 'https://api.github.com/repos/hekmekk/git-team/releases'

local respcode, release = find_release_in_tags(releases_uri, git_team_version)
if respcode == 200 then
  print(string.format("[info] release found for version=%s", git_team_version))
end
if respcode == 404 then
  print(string.format("[info] no release found for version=%s", git_team_version))
  print(string.format('[info] updating git tags', git_team_version))
  interactively_add_git_tag_and_push_to_remote(git_team_version)
  print(string.format('[info] going to create release for version=%s', git_team_version))
  _, release = create_release(releases_uri, github_api_token, git_team_version)
end
if respcode ~= 200 and respcode ~= 404 then
  print(string.format("[error] failure while trying to find release for version=%s", git_team_version))
end

if not release then
  print(string.format('[error] failed to determine upload url for version=%s', git_team_version))
  os.exit(-1)
end

local deb_file_name = ''
for x in git_team_deb_path:gmatch("([^/]+)/?") do deb_file_name = x end

local deb_file = read_file(git_team_deb_path)
if not is_asset_uploaded_already(release['assets'], deb_file_name) then
  print("[debug] asset uploaded already")
  local upload_url_template = release['upload_url']
  local upload_url = string.gsub(upload_url_template, "%{%?name,label%}", string.format('?name=%s', deb_file_name))

  local resp = upload_asset(github_api_token, upload_url, deb_file)
  print(resp)
end

local sha256sum = sha2.sha256hex(deb_file)
if not is_checksum_in_body_already(release, sha256sum) then
  local release_id = release['id']
  local reqbody = json.encode({
    body = string.format('**sha256 checksum:** `%s`', sha256sum)
  })

  local res = update_release(github_api_token, releases_uri, release_id, reqbody)
  print(res)
end

