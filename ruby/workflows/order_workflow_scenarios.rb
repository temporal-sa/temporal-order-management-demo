require_relative 'order_workflow'
require 'temporalio/error'
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
      input = Models::OrderInput.new(input['OrderId'], input['Address']) if input.is_a?(Hash)
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
        logger.error("Failed to charge customer: #{ex.message}")
        compensations.reverse.each do |comp|
          logger.error("Compensating: #{comp}")
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

      #TODO These is not working asynchronously
      activity_handles = order_items.map do |item|
        logger.info("Shipping item: #{item.description}")
        if workflow_type == CHILD
          Temporalio::Workflow::Future.new do
            Temporalio::Workflow.execute_child_workflow(
              Workflows::ShippingChildWorkflow,
              input, item,
              id: "shipment-#{input.order_id}-#{item.id}",
              parent_close_policy: :terminate
            )
          end
        else
          Temporalio::Workflow::Future.new do
            Temporalio::Workflow.execute_activity(
              Activities::ShipOrderActivity,
              input, item,
              start_to_close_timeout: 5,
              retry_policy: @retry_policy
            )
          end
        end
      end

      # Wait for all futures to complete
      results = activity_handles.map(&:wait)

      update_progress("Order Completed", 100, 0)

      tracking_id = Temporalio::Workflow.random.uuid
      Models::OrderOutput.new(tracking_id, input.address).deep_camelize_keys
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
        logger.error("Updating search attributes with order status: #{order_status}")
        key = Temporalio::SearchAttributes::Key.new(
          "OrderStatus",
          Temporalio::SearchAttributes::IndexedValueType::KEYWORD
        )

        update = key.value_set(order_status)

        Temporalio::Workflow.upsert_search_attributes(update)
      end
    end

    workflow_query(name: 'getProgress')
    def query_progress
      @progress
    end

   workflow_signal(name: 'UpdateOrder')
    def update_order(update_input)
      logger.info("Received update order signal with address: #{update_input}")
      @updated_address = update_input['Address']
    end

    # TODO This won't work. Only one of these can be decorated with a specific name, unlike other SDKs where "UpdateOrder" is the name for all three
    # Only one of these will work at a time, and do so by making its name 'UpdateOrder' as that's what the UI calls. Not going to change all the otehr SDKs just for Ruby, which is still in Alpha
    workflow_update(name: 'UpdateOrderUpdate')
    def update_order_update(update_input)
      logger.info("Received update order update with address: #{update_input['Address']}")
      @updated_address = update_input['Address']
      "Updated address: #{@updated_address}"
    end

    workflow_update_validator(:update_order_update)
    def update_order_update_validator(update_input)
      puts("Validating address: #{update_input}")
      unless update_input['Address'][0] =~ /\d/
        logger.info("Rejecting order update, invalid address: #{update_input['Address']}")
        raise Temporalio::Error::ApplicationError.new("Address must start with a digit", type: "invalid-address")
      end
      logger.info("Order update address is valid: #{update_input['Address']}")
    end

    private

    def workflow_type
      Temporalio::Workflow.info.workflow_type
    end
  end
end
