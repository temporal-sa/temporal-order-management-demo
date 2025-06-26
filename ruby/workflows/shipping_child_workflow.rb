require_relative '../activities/order_activities'
require_relative '../shared_objects'

class ShippingChildWorkflow
  attr_accessor :retry_policy

  def initialize
    @retry_policy = nil
  end

  def execute(input, item)
    activities = OrderActivities.new
    activities.ship_order(input, item)
  end
end
