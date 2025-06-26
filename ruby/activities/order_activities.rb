require_relative '../shared_objects'

class OrderActivities
  ERROR_CHARGE_API_UNAVAILABLE = "OrderWorkflowAPIFailure"
  ERROR_INVALID_CREDIT_CARD = "OrderWorkflowNonRecoverableFailure"

  def get_items
    [
      OrderItem.new(id: 654300, description: "Table Top", quantity: 1),
      OrderItem.new(id: 654321, description: "Table Legs", quantity: 2),
      OrderItem.new(id: 654322, description: "Keypad", quantity: 1)
    ]
  end

  def check_fraud(input)
    sleep(1)
    input.order_id
  end

  def prepare_shipment(input)
    sleep(1)
    input.order_id
  end

  def charge_customer(input, type)
    attempt = Temporalio::Activity.info.attempt rescue 1
    error = simulate_external_operation_charge(1, type, attempt)
    case error
    when ERROR_CHARGE_API_UNAVAILABLE
      raise StandardError.new("Charge Customer activity failed, API unavailable")
    when ERROR_INVALID_CREDIT_CARD
      raise StandardError.new("Charge Customer activity failed, card is invalid")
    end
    input.order_id
  end

  def ship_order(input, item)
    sleep(1)
    nil
  end

  def undo_prepare_shipment(input)
    sleep(1)
    input.order_id
  end

  def undo_charge_customer(input)
    sleep(1)
    input.order_id
  end

  private

  def simulate_external_operation(ms)
    sleep(ms / 1000.0)
  end

  def simulate_external_operation_charge(ms, type, attempt)
    simulate_external_operation(ms / attempt)
    attempt < 5 ? type : "NoError"
  end
end
