class OrderOutput
  attr_accessor :tracking_id, :address

  def initialize(tracking_id:, address:)
    @tracking_id = tracking_id
    @address = address
  end
end
