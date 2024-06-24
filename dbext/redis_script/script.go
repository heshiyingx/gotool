package redis_script

import (
	"github.com/redis/go-redis/v9"
)

var GetIncrSeqScript = redis.NewScript(`
local value = redis.call("GET", KEYS[2])
if value  then
	return tonumber(value)
else
	if redis.call("EXISTS", KEYS[1]) ~= 1 then
		return -1
	end

	local seq = redis.call("INCRBY", KEYS[1], 1)
	redis.call("SET", KEYS[2], seq, "EX", 600)
	return seq
end
`)

var IncrExpireScript = redis.NewScript(`
redis.call("SET", KEYS[1], 0,"NX","EX",tonumber(ARGV[1]))
redis.call("EXPIRE",KEYS[1], ARGV[1])
local incrValue =redis.call("INCRBY",KEYS[1], 1)
return incrValue
`)

// DecrZeroDelScript 如果值小于等于0，就删除key
var DecrZeroDelScript = redis.NewScript(`
local value = redis.call("DECR", KEYS[1])
if value <= 0  then
	local status, result = pcall(redis.call,"DEL", KEYS[1])
	if status then
		return 0
	else
		return -1
	end
else
	return value
end
`)

// SetMustGTOldScript 必须要设置比原值大的值。
// 返回 0:比原值小 ，-1操作失败，其他值：设置成功的值
var SetMustGTOldScript = redis.NewScript(`
local value = redis.call("GET", KEYS[1])
local value_num = tonumber(value) or 0 
local ARG_num = tonumber(ARGV[1]) 
if value == false or value == nil then

	local status, result = pcall(redis.call,"SET", KEYS[1], ARGV[1])
	if status then
		return ARGV[1]
	else
		return -1
	end

elseif  ARG_num > value_num  then

	local status, result = pcall(redis.call,"SET", KEYS[1], ARGV[1])
	if status then
		return ARGV[1]
	else
		return -1
	end

else

	return 0

end
`)

// SafeDECRScript key不存在，返回-2，操作失败 -1,其他值是操作成功后，剩下的值
var SafeDECRScript = redis.NewScript(`
local ARG_num = tonumber(ARGV[1]) 
local value = redis.call("GET", KEYS[1])
local value_num = tonumber(value)
if value == false or value == nil then
	return -2
elseif ARG_num <= 0 then
	return -1
elseif value_num >= ARG_num then
	local status, result = pcall(redis.call,"DECRBY", KEYS[1], ARG_num)
	if status then
		return result
	else
		return result
	end
else
	return -1
end

`)
