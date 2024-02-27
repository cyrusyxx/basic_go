local key = KEYS[1]
local cntKey = key .. ":cnt"
local val = ARGV[1]
local ttl = tonumber(redis.call('ttl', key))

-- The key is exist and no expire time
if ttl == -1 then
    return -1
-- The key is not exist or expire time is less than 9min
elseif ttl == -2 or ttl < 540 then
    redis.call('set', key, val)
    redis.call('expire', key, 600)
    redis.call('set', cntKey, 5)
    redis.call('expire', cntKey, 600)
    return 0
-- Send too many times
else
    return -2
end