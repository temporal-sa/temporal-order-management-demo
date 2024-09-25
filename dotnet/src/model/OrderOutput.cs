using System.Text.Json.Serialization;

namespace DotNetOrderManagement.model;

public class OrderOutput
{
    [JsonPropertyName("trackingId")]
    public string TrackingId { get; set; }
    [JsonPropertyName("address")]
    public string Address { get; set; }

    public OrderOutput()
    {
        TrackingId = string.Empty;
        Address = string.Empty;
    }

    public OrderOutput(string trackingId, string address)
    {
        this.TrackingId = trackingId;
        this.Address = address;
    }
}