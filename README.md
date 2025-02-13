1. Goal
The purpose of this document is to provide a detailed and structured overview of the ARM AD Group Creation process. It includes API integration, request workflows, validation rules, and responses. The document ensures:

Automation of AD group onboarding via ARM API.
Defining roles and approval workflow for group creation.
Handling various use cases and validations during AD onboarding.
Providing a structured approach to streamline the process.
2. Background and Strategic Fit
2.1 What is ARM?
ARM (Access Request Management) is the bank’s system for managing AD group access requests.
Historically, ARM did not have an API; instead, requests were manually processed using macro-based spreadsheets.
A new ARM API is now available to automate the process and improve efficiency.
2.2 Why is this needed?
Eliminates manual spreadsheet processing.
Provides a faster and standardized approach for AD onboarding.
Reduces approval delays by automating the workflow.
3. Overview
This document will cover:

API Operations
Approval Workflows
Validation Rules
Successful Requests & Responses
Use Cases and Testing Scenarios
Design Patterns & Architecture
Sequence Diagrams for API Flow
Error Handling & Edge Cases
4. API Operations
The following API endpoints are used for ARM AD Group creation:

Case	Method	API Endpoint	Purpose
1	GET	/arm/request/{requestid}	Get the status of an ARM request (approval state, date, etc.)
2	POST	/arm/request/group	Create a new AD group request
3	POST	/arm/request/group/{group_name}/user	Add a user to an existing AD group
4	DELETE	/arm/request/group/{group_name}	Delete an AD group
5. Design Patterns
The design follows a structured approach:

Microservices-based API architecture.
RESTful API principles for modular & scalable integration.
Role-based validation to ensure security in AD onboarding.
Idempotent API design to prevent duplicate requests.
6. Overall Architecture Diagram
A high-level architecture that includes:

ARM API
AD Group Management
Approval Workflow
Role-based Access Control (RBAC)
(Diagram can be added here)

7. Sequence Diagram
Illustrating the step-by-step API flow for AD group creation and approval.

(Diagram can be added here)

8. Approval Workflow
8.1 General Steps
Request Submission

The requestor submits an ARM request to create an AD group.
Approval Process

Step 1: Associate Manager approval.
Step 2: Selected AIT Manager approval.
Step 3: Primary Approver approval (Final step).
Step 4: Secondary Approver approval (Only if necessary).
Group Creation & Auto-Provisioning

If approvals are completed successfully, the group is automatically provisioned.
9. Use Cases
Use Case 1: AD Object Creation for New App Teams
Scenario:
A new development team needs AD group access. The ARM request includes:

New AD Group Creation (achp_70449_developer_AP)
Primary and Secondary Approvers
User memberships for the group.
ARM Request Payload:

json
Copy
Edit
{
    "RequestUniqueId": "9ae87a...",
    "ExternalRequestorSystem": "Landlord",
    "RequestType": "NEW",
    "RequestorNBK": "ZK71QKG",
    "TicketOnBehalfOf": "On Behalf of Bank Employee",
    "AssociateAndManagerList": [
        {
            "AssociateNBK": "ZKNV4G6",
            "ManagerNBK": "ZKLKRX8"
        }
    ],
    "ARMApplicationCollection": [
        {
            "ApplicationKey": "7684D90F-...",
            "ARMTicketAttributeCollection": [
                {
                    "ARMAttributeType": "AppAttribute",
                    "ARMAttributeKey": "BULK_ACTION",
                    "ARMAttributeValue": "Single Group Request"
                }
            ]
        }
    ]
}
Response:

json
Copy
Edit
{
    "RequestUniqueId": "9ae87a...",
    "SubmitStatus": "Success",
    "ErrorMessage": "",
    "GroupTicketNo": "48743905"
}
Use Case 2: Adding New Members to Existing Groups
A user is added to an existing group using the API.
API Call:

bash
Copy
Edit
curl --request POST \
  --url https://iga-atlas.bankofamerica.com/atlas/api/workOrders/adServiceAccount/modify \
  --header 'Authorization: Bearer ********' \
  --header 'accept: application/json' \
  --header 'content-type: application/json' \
  --data '{
      "justification": "Modify work order testing",
      "initialSourceSystem": {
          "referenceId": "00000001",
          "systemName": "CONTAINER_ACCESS_REQUEST_MANAGEMENT"
      },
      "sourceSystem": {
          "referenceId": "00000002",
          "systemName": "CONTAINER_ACCESS_REQUEST_MANAGEMENT"
      },
      "type": "AD_SERVICE_ID_ADD_GROUPS",
      "serviceAccount": {
          "serviceId": "ZSM9XU5",
          "domain": "CORP",
          "groups": [
              {
                  "groupName": "achp_70449_developer_AP",
                  "domain": "CORP"
              }
          ]
      }
  }'
Response:

json
Copy
Edit
{
    "RequestUniqueId": "9ae87a...",
    "SubmitStatus": "Success",
    "ErrorMessage": "",
    "GroupTicketNo": "48743905"
}
Use Case 3: Checking Status of a Request
API Call:

bash
Copy
Edit
curl --request GET \
  --url https://iga-atlas.bankofamerica.com/atlas/api/processes/<workOrderID> \
  --header 'Authorization: Bearer ********'
Response:

json
Copy
Edit
{
    "ticket_number": "53261908",
    "group_ticket_number": "48743905",
    "request_type": "NEW",
    "status": "Manager Approval(S)",
    "created_time": "11/27/2024 8:19:49 AM"
}
10. Validation Rules
Validation	Description
Primary & Secondary Approver	Must be different persons
Group Name Validation	Auto-appends _AP to group names
ARM Validation	Ensures group does not already exist
User Role Assignment	Ensures proper roles are assigned
Bulk Requests	Not supported in a single API call
11. Testing & QA
Tested API with valid AIT numbers → Passed.
Tested with invalid group names → API rejected.
Checked approval workflow automation → Works as expected.
12. FAQs
1. What happens if a request is rejected?
Rejected requests are closed, and the user must resubmit.

2. Can we create multiple groups in one request?
No, bulk group creation is not supported.

3. How do we check the request status?
Use the GET /arm/request/{requestid} API.
