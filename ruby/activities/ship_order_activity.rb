require 'temporalio/activity'
require_relative '../models/shipping_input'

module Activities
  class ShipOrderActivity < Temporalio::Activity::Definition
    def execute(input)
      context = Temporalio::Activity::Context.current
      context.logger.info("Ship Order activity started, orderId = #{input.order.order_id}, itemId = #{input.item.id}, itemDescription = #{input.item.description}")

      # Item n ships in n*2 seconds, so a fan-out finishes one item at a time; no index means random
      delay_ms = input.index.to_i > 0 ? input.index * 2000 : rand(1000..4000)

      # On a retry, pick up from the last second this item reported
      seconds_done = context.info.heartbeat_details.first || 0
      context.logger.info("Shipping Delay Time: #{delay_ms}, resuming at #{seconds_done}s")

      # Simulate external API call, heartbeating progress each second so a retry can resume
      (seconds_done...(delay_ms / 1000)).each do |second|
        sleep(1)
        context.heartbeat(second + 1)
      end
      sleep((delay_ms % 1000) / 1000.0)
      'SUCCESS'
    end
  end
end
