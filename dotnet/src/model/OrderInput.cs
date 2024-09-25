using System.Text.Json.Serialization;

namespace DotNetOrderManagement.model;

public class OrderInput
{
    [JsonPropertyName("OrderId")]
    public string OrderId { get; set; }
    [JsonPropertyName("Address")]
    public string Address { get; set; }

    public OrderInput()
    {
        OrderId = string.Empty;
        Address = string.Empty;
    }

    public OrderInput(string orderId, string address)
    {
        this.OrderId = orderId;
        this.Address = address;
    }
}