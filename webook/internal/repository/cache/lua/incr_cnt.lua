-- key is the biz key
local key = KEYS[1]
-- cntKey used to choose the field to increment(view, like, or collect)
local cntKey = ARGV[1]
-- delta used to specify +1 or -1
local delta = tonumber(ARGV[2])

-- check if the key exists
local exist=redis.call("EXISTS", key)
if exist == 1 then
    redis.call("HINCRBY", key, cntKey, delta)
    return 1
else
    return 0
end
