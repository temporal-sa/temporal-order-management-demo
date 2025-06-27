require 'temporalio/activity'
require_relative '../shared_objects'

module Activities
  class ShipOrderActivity < Temporalio::Activity::Definition
    def execute(input, item)
      sleep(1)
      'SUCCESS'
    end
  end
end 