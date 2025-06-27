require 'temporalio/activity'
require_relative '../shared_objects'

module Activities
  class UndoPrepareShipmentActivity < Temporalio::Activity::Definition
    def execute(input)
      sleep(1)
      'SUCCESS'
    end
  end
end 