/*
 * Copyright (c) 2013, fromkeith
 * All rights reserved.
 * 
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted provided that the following conditions are met:
 * 
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 * 
 * * Redistributions in binary form must reproduce the above copyright notice, this
 *   list of conditions and the following disclaimer in the documentation and/or
 *   other materials provided with the distribution.
 * 
 * * Neither the name of the fromkeith nor the names of its
 *   contributors may be used to endorse or promote products derived from
 *   this software without specific prior written permission.
 * 
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
 * ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package ec2

import (
    "bytes"
    "net/url"
    "github.com/fromkeith/awsgo"
    "fmt"
    "encoding/xml"
)

const (
    FILTER_architecture = "architecture"
    FILTER_availability_zone = "availability-zone"
    FILTER_block_device_mapping_attach_time = "block-device-mapping.attach-time"
    FILTER_block_device_mapping_delete_on_termination = "block-device-mapping.delete-on-termination"
    FILTER_block_device_mapping_device_name = "block-device-mapping.device-name"
    FILTER_block_device_mapping_status = "block-device-mapping.status"
    FILTER_block_device_mapping_volume_id = "block-device-mapping.volume-id"
    FILTER_client_token = "client-token"
    FILTER_dns_name = "dns-name"
    FILTER_group_id = "group-id"
    FILTER_group_name = "group-name"
    FILTER_image_id = "image-id"
    FILTER_instance_id = "instance-id"
    FILTER_instance_lifecycle = "instance-lifecycle"
    FILTER_instance_state_code = "instance-state-code"
    FILTER_instance_state_name = "instance-state-name"
    FILTER_instance_type = "instance-type"
    FILTER_instance_group_id = "instance.group-id"
    FILTER_instance_group_name = "instance.group-name"
    FILTER_ip_address = "ip-address"
    FILTER_kernel_id = "kernel-id"
    FILTER_key_name = "key-name"
    FILTER_launch_index = "launch-index"
    FILTER_launch_time = "launch-time"
    FILTER_monitoring_state = "monitoring-state"
    FILTER_owner_id = "owner-id"
    FILTER_placement_group_name = "placement-group-name"
    FILTER_platform = "platform"
    FILTER_private_dns_name = "private-dns-name"
    FILTER_private_ip_address = "private-ip-address"
    FILTER_product_code = "product-code"
    FILTER_product_code_type = "product-code.type"
    FILTER_ramdisk_id = "ramdisk-id"
    FILTER_reason = "reason"
    FILTER_requester_id = "requester-id"
    FILTER_reservation_id = "reservation-id"
    FILTER_root_device_name = "root-device-name"
    FILTER_root_device_type = "root-device-type"
    FILTER_source_dest_check = "source-dest-check"
    FILTER_spot_instance_request_id = "spot-instance-request-id"
    FILTER_state_reason_code = "state-reason-code"
    FILTER_state_reason_message = "state-reason-message"
    FILTER_subnet_id = "subnet-id"
    FILTER_tag_key = "tag-key"
    FILTER_tag_value = "tag-value"
    FILTER_virtualization_type = "virtualization-type"
    FILTER_vpc_id = "vpc-id"
    FILTER_hypervisor = "hypervisor"
    FILTER_network_interface_description = "network-interface.description"
    FILTER_network_interface_subnet_id = "network-interface.subnet-id"
    FILTER_network_interface_vpc_id = "network-interface.vpc-id"
    FILTER_network_interface_network_interface_id = "network-interface.network-interface.id"
    FILTER_network_interface_owner_id = "network-interface.owner-id"
    FILTER_network_interface_availability_zone = "network-interface.availability-zone"
    FILTER_network_interface_requester_id = "network-interface.requester-id"
    FILTER_network_interface_requester_managed = "network-interface.requester-managed"
    FILTER_network_interface_status = "network-interface.status"
    FILTER_network_interface_mac_address = "network-interface.mac-address"
    FILTER_network_interface_private_dns_name = "network-interface-private-dns-name"
    FILTER_network_interface_source_destination_check = "network-interface.source-destination-check"
    FILTER_network_interface_group_id = "network-interface.group-id"
    FILTER_network_interface_group_name = "network-interface.group-name"
    FILTER_network_interface_attachment_attachment_id = "network-interface.attachment.attachment-id"
    FILTER_network_interface_attachment_instance_id = "network-interface.attachment.instance-id"
    FILTER_network_interface_attachment_instance_owner_id = "network-interface.attachment.instance-owner-id"
    FILTER_network_interface_addresses_private_ip_address = "network-interface.addresses.private-ip-address"
    FILTER_network_interface_attachment_device_index = "network-interface.attachment.device-index"
    FILTER_network_interface_attachment_status = "network-interface.attachment.status"
    FILTER_network_interface_attachment_attach_time = "network-interface.attachment.attach-time"
    FILTER_network_interface_attachment_delete_on_termination = "network-interface.attachment.delete-on-termination"
    FILTER_network_interface_addresses_primary = "network-interface.addresses.primary"
    FILTER_network_interface_addresses_association_public_ip = "network-interface.addresses.association.public-ip"
    FILTER_network_interface_addresses_association_ip_owner_id = "network-interface.addresses.association.ip-owner-id"
    FILTER_association_public_ip = "association.public-ip"
    FILTER_association_ip_owner_id = "association.ip-owner-id"
    FILTER_association_allocation_id = "association.allocation-id"
    FILTER_association_association_id = "association.association-id"
)

type DescribeFilter struct{
    Name        string
    Value       []string
}

type DescribeInstancesRequest struct {
    awsgo.RequestBuilder

    InstanceIds             []string
    MaxResults              int
    NextToken               string
    Filters                 []DescribeFilter
}

type GroupSet struct {
    GroupId             string  `xml:"groupId"`
    GroupName           string  `xml:"groupName"`
}

type placement struct {
    AvailabilityZone    string  `xml:"availabilityZone"`
    GroupName           string  `xml:"groupName"`
    Tenancy             string  `xml:"tenancy"`
}

type InstanceSet struct {
    InstanceId          string      `xml:"instanceId"`
    ImageId             string      `xml:"imageId"`
    PrivateDnsName      string      `xml:"privateDnsName"`
    DnsName             string      `xml:"dnsName"`
    Reason              string      `xml:"reason"`
    KeyName             string      `xml:"keyName"`
    AmiLaunchIndex      int         `xml:"amiLaunchIndex"`
    ProductCodes        []string        `xml:"productCodes"`
    InstanceType        string      `xml:"instanceType"`
    LaunchTime          string      `xml:"launchTime"`
    Placement           placement   `xml:"placement"`
    Platform            string      `xml:"platform"`
    SubnetId            string      `xml:"subnetId"`
    VpcId               string      `xml:"vpcId"`
    PrivateIpAddress    string      `xml:"privateIpAddress"`
    // add SOOO much more.. arggg
}

type ReservationItem struct {
    ReservationId           string      `xml:"reservationId"`
    OwnerId                 string      `xml:"ownerId"`
    GroupSet                []GroupSet  `xml:"groupSet>item"`
    InstancesSet            []InstanceSet `xml:"instancesSet>item"`

}

type DescribeInstancesResult struct {
    RequestId           string      `xml:"requestId"`
    ReservationSet      []ReservationItem       `xml:"reservationSet>item"`
}


func NewDescribeInstancesRequest() *DescribeInstancesRequest {
    req := new(DescribeInstancesRequest)
    req.Host.Service = "ec2"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Headers = make(map[string]string)
    req.RequestMethod = "GET"
    req.CanonicalUri = "/"
    return req
}

func (gir * DescribeInstancesRequest) VerifyInput() (error) {
    gir.Host.Service = "ec2"
    buf := bytes.Buffer{}
    buf.WriteString(fmt.Sprintf("%s?Action=%s&Version=2014-02-01",
        gir.CanonicalUri,
        url.QueryEscape("DescribeInstances"),
    ))
    for i, inst := range gir.InstanceIds {
        buf.WriteString(fmt.Sprintf("&InstanceId.%d=%s", i + 1, url.QueryEscape(inst)))
    }
    for i, filter := range gir.Filters {
        buf.WriteString(fmt.Sprintf("&Filter.%d.Name=%s", i + 1, url.QueryEscape(filter.Name)))
        for k, val := range filter.Value {
            buf.WriteString(fmt.Sprintf("&Filter.%d.Value.%s=%s", k + 1, url.QueryEscape(val)))
        }
    }
    if gir.MaxResults > 0 {
        buf.WriteString(fmt.Sprintf("&MaxResults=%d", gir.MaxResults))
    }
    if gir.NextToken != "" {
        buf.WriteString(fmt.Sprintf("&NextToken=%s", url.QueryEscape(gir.NextToken)))
    }
    gir.CanonicalUri = buf.String()
    return nil
}

func (gir DescribeInstancesRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    giResponse := new(DescribeInstancesResult)
    fmt.Println(string(response))
    xml.Unmarshal(response, giResponse)
    //json.Unmarshal([]byte(response), giResponse)
    return giResponse
}


func (gir DescribeInstancesRequest) Request() (*DescribeInstancesResult, error) {
    request, err := awsgo.BuildEmptyContentRequest(&gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS2
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*DescribeInstancesResult), err
}
