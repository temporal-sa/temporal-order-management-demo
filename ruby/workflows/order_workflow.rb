require 'temporalio/workflow'
require_relative '../models/order_input'
require_relative '../models/order_output'
require_relative '../models/update_order_input'

module Workflows
  class OrderWorkflow < Temporalio::Workflow::Definition
    workflow_name 'OrderWorkflowHappyPath'
    attr_accessor :progress, :retry_policy

    def initialize
      @progress = 0
      @retry_policy = nil
    end

    def execute(input)
      input = Models::OrderInput.new(input['OrderId'], input['Address']) if input.is_a?(Hash)
      workflow_type = Temporalio::Workflow.info.workflow_type
      logger.info("Order workflow started, type = #{workflow_type}, orderId = #{input.order_id}")

      order_items = Temporalio::Workflow.execute_activity(
        Activities::GetItemsActivity,
        start_to_close_timeout: 5
      )

      Temporalio::Workflow.execute_activity(
        Activities::CheckFraudActivity,
        input,
        start_to_close_timeout: 5,
        retry_policy: @retry_policy
      )
      sleep_fn(1, 25)

      Temporalio::Workflow.execute_activity(
        Activities::PrepareShipmentActivity,
        input,
        start_to_close_timeout: 5,
        retry_policy: @retry_policy
      )
      sleep_fn(1, 50)

      Temporalio::Workflow.execute_activity(
        Activities::ChargeCustomerActivity,
        input, workflow_type,
        start_to_close_timeout: 5,
        retry_policy: @retry_policy
      )
      sleep_fn(3, 75)

      activity_handles = order_items.map do |item|
        logger.info("Shipping item: #{item.description}")
        Temporalio::Workflow::Future.new do
          Temporalio::Workflow.execute_activity(
            Activities::ShipOrderActivity,
            input, item,
            start_to_close_timeout: 5,
            retry_policy: @retry_policy
          )
        end
      end

      # Wait for all futures to complete
      results = activity_handles.map(&:wait)

      sleep_fn(0, 100)

      tracking_id = Temporalio::Workflow.random.uuid
      Models::OrderOutput.new(tracking_id, input.address).deep_camelize_keys
    end

    workflow_query(name: 'getProgress')
    def query_progress
      @progress
    end

    private

    def sleep_fn(seconds, progress)
      sleep(seconds) if seconds > 0
      @progress = progress
    end

    def logger
      @logger ||= Temporalio::Workflow.logger
    end
  end
end
