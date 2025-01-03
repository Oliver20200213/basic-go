-- 从Redis键列表中获取第一个键
local key = KEYS[1]
-- 用户输入的code验证码
local expectedCode = ARGV[1]
local code = redis.call("get",key)
local cntKey = key..":cnt"
-- 转成数字
local cnt = tonumber(redis.call("get",cntKey))

if cnt <= 0 then
    -- 说明用户一直输错,有人在搞事情
    -- 或者已经用过了，也是有人在搞你
    return -1
elseif expectedCode == code then
    -- 输对了
    -- 用完，不能再用了
    redis.call("set",cntKey, -1)
    return 0
else
    -- 输错了
    -- 可验证次数减一
    redis.call("decr",cntKey, -1)
    return -2
end
