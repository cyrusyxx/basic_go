local key = KEYS[1]
local cntKey = key .. ":cnt"

local expectedCode = ARGV[1]

local cnt = tonumber(redis.call('get', cntKey))
local code = redis.call('get', key)

-- cnt is exhausted
if cnt == nil or cnt <= 0 then
    return -1
end

-- code is right
if code == expectedCode then
    redis.call('set', cntKey, 0)
    return 0
-- code is wrong
else
    redis.call('set', cntKey, cnt - 1)
    return -2
end