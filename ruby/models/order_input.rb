class OrderInput
  attr_accessor :order_id, :address

  def initialize(order_id:, address:)
    @order_id = order_id
    @address = address
  end
end
