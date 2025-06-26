require_relative '../activities/order_activities'
require_relative '../shared_objects'
require_relative 'shipping_child_workflow'
require 'securerandom'

class OrderWorkflowScenarios
  BUG = "OrderWorkflowRecoverableFailure"
  CHILD = "OrderWorkflowChildWorkflow"
  SIGNAL = "OrderWorkflowHumanInLoopSignal"
  UPDATE = "OrderWorkflowHumanInLoopUpdate"
  VISIBILITY = "OrderWorkflowAdvancedVisibility"

  attr_accessor :progress, :updated_address, :retry_policy

  def initialize
    @progress = 0
    @updated_address = nil
    @retry_policy = nil
  end

  def execute(args)
    input = args[0]
    workflow_type = "OrderWorkflowScenarios"
    activities = OrderActivities.new
    compensations = []
    order_items = activities.get_items
    update_progress("Check Fraud", 0, 0)
    activities.check_fraud(input)
    update_progress("Prepare Shipment", 25, 1)
    compensations << :undo_prepare_shipment
    activities.prepare_shipment(input)
    update_progress("Charge Customer", 50, 1)
    begin
      compensations << :undo_charge_customer
      activities.charge_customer(input, workflow_type)
    rescue => ex
      compensations.reverse.each do |comp|
        activities.send(comp, input)
      end
      raise ex
    end
    update_progress("Ship Order", 75, 3)
    raise "Simulated bug - fix me!" if workflow_type == BUG
    wait_for_updated_address_or_timeout(input) if [SIGNAL, UPDATE].include?(workflow_type)
    handles = []
    order_items.each do |item|
      handles << Thread.new { ship_item_async(input, item, workflow_type) }
    end
    handles.each(&:join)
    update_progress("Order Completed", 100, 0)
    tracking_id = SecureRandom.uuid
    OrderOutput.new(tracking_id: tracking_id, address: input.address)
  end

  def ship_item_async(input, item, workflow_type)
    if workflow_type == CHILD
      ShippingChildWorkflow.new.execute(input, item)
    else
      OrderActivities.new.ship_order(input, item)
    end
  end

  def wait_for_updated_address_or_timeout(input)
    # Simulate waiting for up to 60 seconds for updated address
    60.times do
      break if @updated_address
      sleep(1)
    end
    input.address = @updated_address if @updated_address
  end

  def update_progress(order_status, progress, sleep_time)
    @progress = progress
    sleep(sleep_time) if sleep_time > 0
  end

  def query_progress
    @progress
  end

  def update_order_signal(update_input)
    @updated_address = update_input.address
  end

  def update_order_update(update_input)
    @updated_address = update_input.address
    "Updated address: #{@updated_address}"
  end

  def update_order_validator(update_input)
    raise "Address must start with a digit" unless update_input.address[0] =~ /\d/
  end
end
