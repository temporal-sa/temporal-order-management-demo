require 'temporalio/activity'
require_relative '../shared_objects'

module Activities
  class GetItemsActivity < Temporalio::Activity::Definition
    def execute
      [
        OrderItem.new(654300, "Table Top", 1),
        OrderItem.new(654321, "Table Legs", 2),
        OrderItem.new(654322, "Keypad", 1)
      ]
    end
  end
end 