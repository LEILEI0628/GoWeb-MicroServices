local key = KEYS[1]
-- 对应hincrby中的field
local cntKey = ARGV[1]
-- +1 者 -1
local delta = tonumber(ARGV[2])
local exists = redis.call("EXISTS", key)
if exists == 1 then
    redis.call("HINCRBY", key, cntKey, delta)
    -- 自增成功
    return 1
else
    -- 自增失败
    return 0
end