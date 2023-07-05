require "redis"
require "benchmark"

# Khởi tạo Redis connection
redis = Redis.new

# Benchmark không pipeline
time = Benchmark.realtime do
  (0...10000).each do |i|
    redis.set("key#{i}", "value#{i}")
    redis.get("key#{i}")
  end
end
puts "Elapsed time without pipeline: #{(time * 1000).round(2)}ms"

# Benchmark pipeline
time = Benchmark.realtime do
  redis.pipelined do
    (0...10000).each do |i|
      redis.set("key#{i}", "value#{i}")
      redis.get("key#{i}")
    end
  end
end
puts "Elapsed time with pipeline: #{(time * 1000).round(2)}ms"

# Đóng Redis connection
redis.close
