require_relative '../activities/order_activities'
require_relative '../shared_objects'
require 'securerandom'

class OrderWorkflow
  attr_accessor :progress, :retry_policy

  def initialize
    @progress = 0
    @retry_policy = nil
  end

  def execute(input)
    workflow_type = "OrderWorkflowHappyPath"
    @retry_policy = nil
    activities = OrderActivities.new
    order_items = activities.get_items
    activities.check_fraud(input)
    sleep_fn(1, 25)
    activities.prepare_shipment(input)
    sleep_fn(1, 50)
    activities.charge_customer(input, workflow_type)
    sleep_fn(3, 75)
    handles = []
    order_items.each do |item|
      handles << Thread.new { activities.ship_order(input, item) }
    end
    handles.each(&:join)
    sleep_fn(0, 100)
    tracking_id = SecureRandom.uuid
    OrderOutput.new(tracking_id: tracking_id, address: input.address)
  end

  def query_progress
    @progress
  end

  def sleep_fn(seconds, progress)
    sleep(seconds) if seconds > 0
    @progress = progress
  end
end
