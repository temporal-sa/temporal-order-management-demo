class OrderItem
  attr_accessor :id, :description, :quantity

  def initialize(id:, description:, quantity:)
    @id = id
    @description = description
    @quantity = quantity
  end
end
