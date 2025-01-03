-- 你的验证码在Redis上的key
-- KEYS[1]第一个参数是c.key(biz,phone)
-- phone_code:login:155xxxxxxxx
local key = KEYS[1]
-- 验证次数，一个验证码，最多重复验证三次，这个记录了还可以验证几次
-- phone_code:login:155xxxxxxxx:cnt
local cntKey = key .. ":cnt"
-- ARGV[1]也就是上面Eval中的code，验证码，
-- 你的验证码：例如12345
local val = ARGV[1]
-- 过期时间,   redis.call("ttl", key)获取key的过期时间
-- tonumber转成数字
-- 执行redis的ttl命令获取key的过期时间
local ttl = tonumber(redis.call("ttl", key))
-- ttl的值-1 -2
-- -1表示key存在但是没有过期时间，也就是有人手动设置了这个key，但是没给过期时间
-- -2表示key不存在
if ttl == -1 then
    -- reids ttl 值为 -1,key存在，但是没有过期时间
    return -2 --这里的-2和ttl的-2不一样，这里的-2表示系统错误，自己定义的
elseif ttl == -2 or ttl < 540 then
    -- 540 == 600-60=9分钟  可以发短信
    redis.call("set", key, val)  --key为键 val为值
    redis.call("expire", key, 600) --设置key的过期时间为600s
    redis.call("set", cntKey, 3)    --设置cntKey为3
    redis.call("expire", cntKey, 600)  --cntKey的过期时间为600s
    -- return 0 表示一切正常
    return 0
else
    -- 自己定义-1 表示发送太频繁
    return -1
end



