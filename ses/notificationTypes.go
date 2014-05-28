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

package ses

// see http://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html#bounce-types
const (
    NotificationType_Complaint = "Complaint"
    NotificationType_Bounce = "Bounce"

    BounceType_Undetermined = "Undetermined"
    SubBounceType_Undetermined = "Undetermined"
    // Permanent Bounces
    BounceType_Permanent = "Permanent"
    SubBounceType_General = "General"
    SubBounceType_NoEmail = "NoEmail"
    SubBounceType_Suppressed = "Suppressed"
    // Transient Bounces
    BounceType_Transient = "Transient"
    // also valid for transient : SubBounceType_General = "General"
    SubBounceType_MailboxFull = "MailboxFull"
    SubBounceType_MessageTooLarge = "MessageTooLarge"
    SubBounceType_ContentRejected = "ContentRejected"
    SubBounceType_AttachmentRejected = "AttachmentRejected"

    // Complaint types
    ComplaintType_Abuse = "abuse"
    ComplaintType_AuthFailure = "auth-failure"
    ComplaintType_Fraud = "fraud"
    ComplaintType_NotSpam = "not-spam"
    ComplaintType_Other = "other"
    ComplaintType_Virus = "virus"
)


type BouncedRecipient struct {
    Status              string `json:"status"`
    Action              string `json:"action"`
    DiagnosticCode      string `json:"diagnosticCode"`
    EmailAddress        string `json:"emailAddress"`
}

type SesBounce struct {
    BounceSubType           string      `json:"bounceSubType"`
    BounceType              string      `json:"bounceType"`
    ReportingMTA            string      `json:"reportingMTA"`
    BouncedRecipients       []BouncedRecipient `json:"bouncedRecipients"`
    Timestamp               string      `JSON:"timestamp"`
    FeedbackId              string      `JSON:"feedbackId"`
}
type SesMail struct {
    Timestamp               string      `json:"timestamp"`
    Source                  string      `json:"source"`
    MessageId               string      `json:"messageId"`
    Destination             []string      `json:"destination"`
}

type ComplainedRecipient  struct {
    EmailAddress            string      `json:"emailAddress"`
}
type SesComplaint struct {
    UserAgent               string      `json:"userAgent"`
    ComplaintFeedbackType   string      `json:"complaintFeedbackType"`
    ComplainedRecipients    []ComplainedRecipient `json:"complainedRecipients"`
    Timestamp               string      `JSON:"timestamp"`
    FeedbackId              string      `JSON:"feedbackId"`
}

/*
    Ses notifications as defined: http://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-examples.html
    You can get these sent to you via SNS.
*/
type SesFeedbackNotification struct {
    NotificationType        string      `json:"notificationType"`
    Bounce                  *SesBounce  `json:"bounce,omitempty"`
    Mail                    *SesMail    `json:"mail,omitempty"`
    Complaint               *SesComplaint `json:"complaint,omitempty"`
}