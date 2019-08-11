--[[ Dependencies
sudo apt install luarocks libssl-dev
sudo luarocks install argparse
sudo luarocks install luasec
sudo luarocks install inspect
sudo luarocks install lua-cjson
sudo luarocks install sha2
--]]

--[[ Purge a release (and respective tag)
curl -I -f -H 'Authorization: token GITHUB_API_TOKEN' -X DELETE https://api.github.com/repos/hekmekk/git-team/releases/:id
git tag -d VERSION
git push origin --delete VERSION
--]]

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

  if respbody and respbody ~= "" then
    return respcode, json.decode(respbody)
  end

  return respcode, nil
end

function delete_asset(github_api_token, releases_uri, asset_id)
  local respbody = {}
  local result, respcode, respheaders, respstatus = https.request {
      method = 'DELETE',
      url = releases_uri .. '/assets/' .. asset_id,
      headers = {
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

function update_release(github_api_token, releases_uri, release_id, request_body)
  local respbody = {} -- for the response body

  local result, respcode, respheaders, respstatus = https.request {
      method = 'PATCH',
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
  local reqbody = json.encode({ tag_name = version, body = '', name = version })
  local respbody = {}

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

function get_current_checksum(release)
  if release and release['body'] then
    return release['body']:match("%*%*sha256 checksum:%*%* `(%w+)`")
  end
  return nil
end

function find_asset_by_name(release, asset_name)
  if not release or not release['assets'] then
    return nil
  end
  for _, asset in pairs(release['assets']) do
    if asset['name'] and asset['name'] == asset_name then
      return asset
    end
  end
  return nil
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
    os.execute('git push origin --tags')
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
  print(string.format("[info ] release with id=%s found for version=%s", release['id'], git_team_version))
end
if respcode == 404 then
  print(string.format("[info ] no release found for version=%s", git_team_version))
  interactively_add_git_tag_and_push_to_remote(git_team_version)
  print(string.format('[info ] latest commit has been tagged', version))
  print(string.format('[info ] remote tags have been updated', version))
  ret, release = create_release(releases_uri, github_api_token, git_team_version) -- note: MUST NOT use local here... or else 'release' will not be visible after 'end'
  if ret == 201 then
    print(string.format('[info ] release with id=%s has been created successfully', release['id']))
  else
    print(string.format('[error] failed to create release for version %s', git_team_version))
  end
end
if respcode ~= 200 and respcode ~= 404 then
  print(string.format('[error] failure while trying to find release for version=%s', git_team_version))
end

if not release then
  print(string.format('[error] failed to determine upload url for version=%s', git_team_version))
  os.exit(-1)
end

local deb_file_name = ''
for x in git_team_deb_path:gmatch("([^/]+)/?") do deb_file_name = x end

local local_deb_file = read_file(git_team_deb_path)
local remote_checksum = get_current_checksum(release)
local remote_asset = find_asset_by_name(release, deb_file_name)

if remote_asset and remote_checksum then
  local browser_download_url = remote_asset['browser_download_url']
  if browser_download_url then
    -- TODO: fix this logic - download file and compute checksum to compare against the current remote; exit if equal
    -- local remote_deb_file = ltn12.sink.file(io.open('/tmp/git-team-remote.deb', 'w'))
    -- local result, code, respheaders, respstatus = https.request {
      -- url = browser_download_url,
      -- sink = remote_deb_file
    -- }

    print('[info ] assuming consistency regarding remote asset and checksum - done.')
    os.exit(0)
  end
end

-- checksum
print('[info ] updating checksum ...')
local release_id = release['id']
local reqbody = json.encode({
  body = string.format('**sha256 checksum:** `%s`', sha2.sha256hex(local_deb_file))
})

local respcode2, body2 = update_release(github_api_token, releases_uri, release_id, reqbody)

if respcode2 == 200 then
  print("[info ] checksum added successfully")
else
  print("[error] failed to update release with checksum")
  print(string.format("[debug] release_id=%s", release_id))
  print(string.format("[debug] code=%s", respcode2))
  print(string.format("[debug] body=%s", inspect(body2)))
  os.exit(-1)
end

-- asset
if remote_asset then
  print('[info ] deleting the current asset ...')
  local asset_id = remote_asset['id']
  local code, res = delete_asset(github_api_token, releases_uri, asset_id)
  if code == 204 then
    print('[info ] asset deleted successfully')
  else
    print(string.format('[error] failed to delete asset with id=%s', asset_id))
    print(string.format("[debug] code=%s", code))
    print(string.format("[debug] body=%s", inspect(res)))
    os.exit(-1)
  end
end

print('[info ] creating a new asset ...')
local upload_url_template = release['upload_url']
local upload_url = string.gsub(upload_url_template, "%{%?name,label%}", string.format('?name=%s', deb_file_name))

local respcode, body = upload_asset(github_api_token, upload_url, local_deb_file)
if respcode == 201 then
  print("[info ] asset created successfully")
else
  print("[error] failed to upload asset")
  print(string.format("[debug] code=%s", respcode))
  print(string.format("[debug] body=%s", inspect(body)))
end

