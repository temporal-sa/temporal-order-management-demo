require 'temporalio/workflow'
require 'temporalio/error'
require_relative '../shared_objects'
require_relative 'shipping_child_workflow'

module Workflows
  class OrderWorkflowScenarios < OrderWorkflow
    workflow_dynamic

    BUG = "OrderWorkflowRecoverableFailure"
    CHILD = "OrderWorkflowChildWorkflow"
    SIGNAL = "OrderWorkflowHumanInLoopSignal"
    UPDATE = "OrderWorkflowHumanInLoopUpdate"
    VISIBILITY = "OrderWorkflowAdvancedVisibility"

    attr_accessor :progress, :updated_address, :retry_policy

    def initialize
      super
      @updated_address = nil
    end

    def execute(input)
      input = OrderInput.new(input['OrderId'], input['Address']) if input.is_a?(Hash)
      workflow_type = Temporalio::Workflow.info.workflow_type
      logger.info("Dynamic Order workflow started, type = #{workflow_type}, orderId = #{input.order_id}")

      compensations = []

      order_items = Temporalio::Workflow.execute_activity(
        Activities::GetItemsActivity,
        start_to_close_timeout: 5
      )

      update_progress("Check Fraud", 0, 0)

      Temporalio::Workflow.execute_activity(
        Activities::CheckFraudActivity,
        input,
        start_to_close_timeout: 5,
        retry_policy: @retry_policy
      )

      update_progress("Prepare Shipment", 25, 1)

      compensations << Activities::UndoPrepareShipmentActivity
      Temporalio::Workflow.execute_activity(
        Activities::PrepareShipmentActivity,
        input,
        start_to_close_timeout: 5,
        retry_policy: @retry_policy
      )

      update_progress("Charge Customer", 50, 1)

      begin
        compensations << Activities::UndoChargeCustomerActivity
        Temporalio::Workflow.execute_activity(
          Activities::ChargeCustomerActivity,
          input, workflow_type,
          start_to_close_timeout: 5,
          retry_policy: @retry_policy
        )
      rescue StandardError => ex
        logger.error("Failed to charge customer", ex)
        compensations.reverse.each do |comp|
          Temporalio::Workflow.execute_activity(
            comp,
            input,
            start_to_close_timeout: 10,
            retry_policy: @retry_policy
          )
        end
        raise ex
      end

      update_progress("Ship Order", 75, 3)

      raise "Simulated bug - fix me!" if workflow_type == BUG

      wait_for_updated_address_or_timeout(input) if [SIGNAL, UPDATE].include?(workflow_type)

      handles = []
      order_items.each do |item|
        logger.info("Shipping item: #{item.description}")
        ship_item_async(input, item, workflow_type)
      end

      update_progress("Order Completed", 100, 0)

      tracking_id = Temporalio::Workflow.random.uuid
      OrderOutput.new(tracking_id, input.address)
    end

    def ship_item_async(input, item, workflow_type)
      if workflow_type == CHILD
        Temporalio::Workflow.execute_child_workflow(
          Workflows::ShippingChildWorkflow,
          input, item,
          id: "shipment-#{input.order_id}-#{item.id}",
          parent_close_policy: :terminate
        )
      else
        Temporalio::Workflow.execute_activity(
          Activities::ShipOrderActivity,
          input, item,
          start_to_close_timeout: 5,
          retry_policy: @retry_policy
        )
      end
    end

    def wait_for_updated_address_or_timeout(input)
      logger.info("Waiting up to 60 seconds for updated address")
      begin
        Temporalio::Workflow.timeout(60) { Temporalio::Workflow.wait_condition { @updated_address } }
        input.address = @updated_address
      rescue Timeout::Error
        logger.info("Updated address was not received within 60 seconds")
      end
    end

    def update_progress(order_status, progress, sleep_time)
      @progress = progress
      sleep(sleep_time) if sleep_time > 0
      if workflow_type == VISIBILITY
        Temporalio::Workflow.upsert_search_attributes('OrderStatus' => order_status)
      end
    end

    workflow_query(name: 'getProgress')
    def query_progress
      @progress
    end

    workflow_signal(name: 'UpdateOrderSignal')
    def update_order_signal(update_input)
      logger.info("Received update order signal with address: #{update_input.address}")
      @updated_address = update_input.address
    end

    workflow_update(name: 'UpdateOrder')
    def update_order_update(update_input)
      logger.info("Received update order update with address: #{update_input.address}")
      @updated_address = update_input.address
      "Updated address: #{@updated_address}"
    end

    def update_order_validator(update_input)
      unless update_input.address[0] =~ /\d/
        logger.info("Rejecting order update, invalid address: #{update_input.address}")
        raise Temporalio::Error::ApplicationError.new("Address must start with a digit", type: "invalid-address")
      end
      logger.info("Order update address is valid: #{update_input.address}")
    end

    private

    def workflow_type
      Temporalio::Workflow.info.workflow_type
    end
  end
end
