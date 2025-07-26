require 'temporalio/activity'

module Activities
  class ShipOrderActivity < Temporalio::Activity::Definition
    def execute(input, item)
      sleep(1)
      'SUCCESS'
    end
  end
end
