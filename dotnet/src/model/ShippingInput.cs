using System.Text.Json.Serialization;

namespace DotNetOrderManagement.model;

public class ShippingInput
{
    [JsonPropertyName("Order")]
    public OrderInput Order { get; set; }
    [JsonPropertyName("Item")]
    public OrderItem Item { get; set; }
    // item n ships in n*2 seconds; 0 means a random delay
    [JsonPropertyName("Index")]
    public int Index { get; set; }

    public ShippingInput()
    {
        Order = new OrderInput();
        Item = new OrderItem();
        Index = 0;
    }

    public ShippingInput(OrderInput order, OrderItem item) : this(order, item, 0)
    {
    }

    public ShippingInput(OrderInput order, OrderItem item, int index)
    {
        this.Order = order;
        this.Item = item;
        this.Index = index;
    }
}
