require 'temporalio/workflow'
require_relative '../shared_objects'

module Workflows
  class ShippingChildWorkflow < Temporalio::Workflow::Definition
    attr_accessor :retry_policy

    def initialize
      @retry_policy = nil
    end

    def execute(input, item)
      logger.info("Shipping workflow started, orderId = #{input.order_id}")

      Temporalio::Workflow.execute_activity(
        Activities::ShipOrderActivity,
        input, item,
        start_to_close_timeout: 5,
        retry_policy: @retry_policy
      )
    end

    private

    def logger
      @logger ||= Temporalio::Workflow.logger
    end
  end
end
