# go-internxt-drive

A pure Go library allowing you to access the [internxt drive API](https://api.internxt.com/drive/).

Please refer to `docs/login.md` for an example on how to get started.

## Status

Implementation is **WIP** and here's the current status.

### Auth

| Method | Endpoint                            | Description                                   | Implemented |
| ------ | ----------------------------------- | --------------------------------------------- | ----------- |
| POST   | `/drive/auth/login`                 | Get security details to log in                | Yes         |
| POST   | `/drive/auth/login/access`          | Access user account                           | Yes         |
| GET    | `/api/auth/logout`                  | Log out of the account                        | No          |
| GET    | `/api/auth/tfa`                     | Get two-factor authentication                 | No          |
| PUT    | `/api/auth/tfa`                     | Update two-factor authentication              | No          |
| DELETE | `/api/auth/tfa`                     | Delete two-factor authentication              | No          |
| GET    | `/api/auth/are-credentials-correct` | Check if current user credentials are correct | No          |

### File

| Method | Endpoint                         | Description                           | Implemented |
| ------ | -------------------------------- | ------------------------------------- | ----------- |
| POST   | `/api/files`                     | Create File                           | No          |
| GET    | `/api/files`                     | —                                     | No          |
| GET    | `/api/files/count`               | —                                     | No          |
| GET    | `/api/files/{uuid}/meta`         | —                                     | No          |
| PUT    | `/api/files/{uuid}/meta`         | Update File data                      | No          |
| PUT    | `/api/files/{uuid}`              | —                                     | No          |
| PATCH  | `/api/files/{uuid}`              | —                                     | No          |
| DELETE | `/api/files/{uuid}`              | Delete file from storage and database | No          |
| GET    | `/api/files/recents`             | —                                     | No          |
| GET    | `/api/files/meta`                | —                                     | No          |
| POST   | `/api/files/thumbnail`           | Create Thumbnail                      | No          |
| DELETE | `/api/files/{bucketId}/{fileId}` | Delete file from storage by fileId    | No          |

### Share

| Method | Endpoint                               | Description                        | Implemented |
| ------ | -------------------------------------- | ---------------------------------- | ----------- |
| GET    | `/api/storage/share/domains`           | —                                  | No          |
| GET    | `/api/storage/share/list`              | Get share list                     | No          |
| GET    | `/api/storage/share/{token}`           | Get share by token                 | No          |
| PUT    | `/api/storage/share/{token}/view`      | Increment share view by token      | No          |
| PUT    | `/api/storage/share/{shareId}`         | Update share by id                 | No          |
| DELETE | `/api/storage/share/{shareId}`         | Delete share by id                 | No          |
| POST   | `/api/storage/share/file/{fileId}`     | Create share for file              | No          |
| POST   | `/api/storage/share/folder/{folderId}` | Create share for folder            | No          |
| GET    | `/api/storage/share/down/files`        | Get all files by token paginated   | No          |
| GET    | `/api/storage/share/down/folders`      | Get all folders by token paginated | No          |

### Folder

| Method | Endpoint                                        | Description                     | Implemented |
| ------ | ----------------------------------------------- | ------------------------------- | ----------- |
| POST   | `/api/folders`                                  | Create Folder                   | No          |
| DELETE | `/api/folders`                                  | —                               | No          |
| GET    | `/api/folders`                                  | —                               | No          |
| GET    | `/api/folders/count`                            | —                               | No          |
| GET    | `/api/folders/content/{uuid}/files`             | —                               | No          |
| GET    | `/api/folders/{id}/files`                       | —                               | No          |
| GET    | `/api/folders/{id}/file`                        | —                               | No          |
| GET    | `/api/folders/content/{uuid}/folders`           | —                               | No          |
| GET    | `/api/folders/content/{uuid}/folders/existence` | Checks folders existence (GET)  | No          |
| POST   | `/api/folders/content/{uuid}/folders/existence` | Checks folders existence (POST) | No          |
| POST   | `/api/folders/content/{uuid}/files/existence`   | Checks files existence          | No          |
| GET    | `/api/folders/content/{uuid}`                   | Gets folder content             | No          |
| GET    | `/api/folders/{id}/folders`                     | —                               | No          |
| GET    | `/api/folders/{uuid}/meta`                      | —                               | No          |
| PUT    | `/api/folders/{uuid}/meta`                      | —                               | No          |
| GET    | `/api/folders/{uuid}/ancestors`                 | —                               | No          |
| GET    | `/api/folders/{uuid}/tree`                      | —                               | No          |
| GET    | `/api/folders/{id}/metadata`                    | —                               | No          |
| GET    | `/api/folders/{uuid}/size`                      | —                               | No          |
| PATCH  | `/api/folders/{uuid}`                           | —                               | No          |
| DELETE | `/api/folders/{uuid}`                           | Delete Folder                   | No          |
| GET    | `/api/folders/meta`                             | —                               | No          |

## Sharing

| Method | Endpoint                                              | Description                          | Implemented |
| ------ | ----------------------------------------------------- | ------------------------------------ | ----------- |
| GET    | `/api/sharings/{sharingId}/meta`                      | Get sharing metadata                 | No          |
| GET    | `/api/sharings/public/{sharingId}/item`               | Get sharing item info                | No          |
| PATCH  | `/api/sharings/{sharingId}/password`                  | Set password for public sharing      | No          |
| DELETE | `/api/sharings/{sharingId}/password`                  | Remove password from public sharing  | No          |
| GET    | `/api/sharings/{itemType}/{itemId}/invites`           | List invites for an item             | No          |
| PUT    | `/api/sharings/{itemType}/{itemId}/type`              | Change sharing type for an item      | No          |
| GET    | `/api/sharings/{itemType}/{itemId}/type`              | Get sharing type for an item         | No          |
| GET    | `/api/sharings/{itemType}/{itemId}/info`              | Get info related to item sharing     | No          |
| GET    | `/api/sharings/invites`                               | Get all invites received by the user | No          |
| POST   | `/api/sharings/invites/send`                          | Send a sharing invite                | No          |
| GET    | `/api/sharings/invites/{id}/validate`                 | Validate a sharing invite            | No          |
| POST   | `/api/sharings/invites/{id}/accept`                   | Accept a sharing invite              | No          |
| DELETE | `/api/sharings/invites/{id}`                          | Delete a sharing invite              | No          |
| GET    | `/api/sharings/items/{sharedFolderId}/folders`        | Get folders in a shared folder       | No          |
| GET    | `/api/sharings/items/{sharedFolderId}/files`          | Get files in a shared folder         | No          |
| GET    | `/api/sharings/public/items/{sharedFolderId}/files`   | Get files in a public share          | No          |
| GET    | `/api/sharings/public/items/{sharedFolderId}/folders` | Get folders in a public share        | No          |
| POST   | `/api/sharings`                                       | Share an item                        | No          |
| DELETE | `/api/sharings/{itemType}/{itemId}`                   | Stop sharing an item                 | No          |
| GET    | `/api/sharings/roles`                                 | List sharing roles                   | No          |
| GET    | `/api/sharings/{sharingId}/role`                      | Get role of a sharing                | No          |
| PUT    | `/api/sharings/{sharingId}/role`                      | Set role of a sharing                | No          |
| DELETE | `/api/sharings/{sharingId}/roles/{sharingRoleId}`     | Delete a sharing role                | No          |
| GET    | `/api/sharings/shared-with-me/folders`                | Folders shared with me               | No          |
| GET    | `/api/sharings/shared-by-me/folders`                  | Folders I’ve shared                  | No          |
| GET    | `/api/sharings/folders`                               | All folder sharings                  | No          |
| GET    | `/api/sharings/files`                                 | All file sharings                    | No          |
| GET    | `/api/sharings/shared-with/{itemType}/{itemId}`       | Users with access to an item         | No          |
| GET    | `/api/sharings/shared-with/{folderId}`                | Users with access to a folder        | No          |
| DELETE | `/api/sharings/{itemType}/{itemId}/users/{userId}`    | Remove a user from a shared item     | No          |
| GET    | `/api/sharings/public/{id}/folder/size`               | Get size of a public shared folder   | No          |

### Workspaces

| Method | Endpoint                                                                       | Description                       | Implemented |
| ------ | ------------------------------------------------------------------------------ | --------------------------------- | ----------- |
| GET    | `/api/workspaces`                                                              | Get available workspaces          | No          |
| GET    | `/api/workspaces/pending-setup`                                                | Get workspaces pending setup      | No          |
| GET    | `/api/workspaces/invitations`                                                  | Get pending workspace invitations | No          |
| POST   | `/api/workspaces/invitations/accept`                                           | Accept a workspace invitation     | No          |
| GET    | `/api/workspaces/invitations/{inviteId}/validate`                              | Validate a workspace invitation   | No          |
| DELETE | `/api/workspaces/invitations/{inviteId}`                                       | Decline a workspace invitation    | No          |
| GET    | `/api/workspaces/teams/{teamId}/members`                                       | List team members                 | No          |
| PATCH  | `/api/workspaces/teams/{teamId}`                                               | Edit team data                    | No          |
| DELETE | `/api/workspaces/{workspaceId}/teams/{teamId}`                                 | Delete a team                     | No          |
| POST   | `/api/workspaces/teams/{teamId}/user/{userUuid}`                               | Add user to team                  | No          |
| DELETE | `/api/workspaces/teams/{teamId}/user/{userUuid}`                               | Remove user from team             | No          |
| GET    | `/api/workspaces/{workspaceId}/files`                                          | List workspace files              | No          |
| POST   | `/api/workspaces/{workspaceId}/files`                                          | Create workspace file             | No          |
| GET    | `/api/workspaces/{workspaceId}/folders`                                        | List workspace folders            | No          |
| POST   | `/api/workspaces/{workspaceId}/folders`                                        | Create workspace folder           | No          |
| PATCH  | `/api/workspaces/{workspaceId}/teams/{teamId}/manager`                         | Change team manager               | No          |
| GET    | `/api/workspaces/{workspaceId}/invitations`                                    | List workspace invitations        | No          |
| PATCH  | `/api/workspaces/{workspaceId}/setup`                                          | Setup an initialized workspace    | No          |
| POST   | `/api/workspaces/{workspaceId}/avatar`                                         | Upload workspace avatar           | No          |
| DELETE | `/api/workspaces/{workspaceId}/avatar`                                         | Delete workspace avatar           | No          |
| GET    | `/api/workspaces/{workspaceId}/credentials`                                    | Get workspace credentials         | No          |
| GET    | `/api/workspaces/{workspaceId}/usage`                                          | Get workspace usage               | No          |
| PATCH  | `/api/workspaces/{workspaceId}/members/{memberId}/usage`                       | Change member assigned space      | No          |
| GET    | `/api/workspaces/{workspaceId}/members`                                        | List workspace members            | No          |
| POST   | `/api/workspaces/{workspaceId}/members/invite`                                 | Invite user to workspace          | No          |
| POST   | `/api/workspaces/{workspaceId}/teams`                                          | Create workspace team             | No          |
| GET    | `/api/workspaces/{workspaceId}/teams`                                          | List workspace teams              | No          |
| GET    | `/api/workspaces/{workspaceId}/usage/member`                                   | Get user usage in drive           | No          |
| POST   | `/api/workspaces/{workspaceId}/shared`                                         | Share item to workspace           | No          |
| GET    | `/api/workspaces/{workspaceId}/teams/{teamId}/shared/files`                    | List shared files in team         | No          |
| GET    | `/api/workspaces/{workspaceId}/shared/files`                                   | List shared files                 | No          |
| GET    | `/api/workspaces/{workspaceId}/teams/{teamId}/shared/folders`                  | List shared folders in team       | No          |
| GET    | `/api/workspaces/{workspaceId}/shared/folders`                                 | List shared folders               | No          |
| GET    | `/api/workspaces/{workspaceId}/teams/{teamId}/shared/{sharedFolderId}/folders` | List folders in a shared folder   | No          |
| GET    | `/api/workspaces/{workspaceId}/shared/{sharedFolderId}/folders`                | List folders in a shared folder   | No          |
| GET    | `/api/workspaces/{workspaceId}/teams/{teamId}/shared/{sharedFolderId}/files`   | List files in a shared folder     | No          |
| GET    | `/api/workspaces/{workspaceId}/shared/{sharedFolderId}/files`                  | List files in a shared folder     | No          |
| GET    | `/api/workspaces/{workspaceId}/shared/{itemType}/{itemId}/shared-with`         | List shares for an item           | No          |
| GET    | `/api/workspaces/{workspaceId}/trash`                                          | Get workspace trash               | No          |
| DELETE | `/api/workspaces/{workspaceId}/trash`                                          | Empty workspace trash             | No          |
| GET    | `/api/workspaces/{workspaceId}/folders/{folderUuid}/folders`                   | List subfolders                   | No          |
| GET    | `/api/workspaces/{workspaceId}/folders/{folderUuid}/files`                     | List files in a folder            | No          |
| PATCH  | `/api/workspaces/{workspaceId}/teams/{teamId}/members/{memberId}/role`         | Change member’s role              | No          |
| PATCH  | `/api/workspaces/{workspaceId}`                                                | Edit workspace details            | No          |
| GET    | `/api/workspaces/{workspaceId}`                                                | Get workspace details             | No          |
| DELETE | `/api/workspaces/{workspaceId}/members/leave`                                  | Leave workspace                   | No          |
| GET    | `/api/workspaces/{workspaceId}/members/{memberId}`                             | Get member details                | No          |
| DELETE | `/api/workspaces/{workspaceId}/members/{memberId}`                             | Remove member                     | No          |
| PATCH  | `/api/workspaces/{workspaceId}/members/{memberId}/deactivate`                  | Deactivate a member               | No          |
| PATCH  | `/api/workspaces/{workspaceId}/members/{memberId}/activate`                    | Activate a member                 | No          |

### User

| Method | Endpoint                                                 | Description                         | Implemented |
| ------ | -------------------------------------------------------- | ----------------------------------- | ----------- |
| POST   | `/api/users`                                             | Create a user                       | No          |
| GET    | `/api/users/user/{email}`                                | Get user data & subscription status | No          |
| POST   | `/api/users/pre-created-users/register`                  | Register a pre-created user         | No          |
| POST   | `/api/users/pre-create`                                  | Pre-create a user                   | No          |
| GET    | `/api/users/c/{uuid}`                                    | Get user credentials                | No          |
| GET    | `/api/users/refresh`                                     | Refresh session token               | No          |
| PATCH  | `/api/users/password`                                    | Change user password                | No          |
| POST   | `/api/users/recover-account`                             | Request account recovery            | No          |
| PUT    | `/api/users/recover-account`                             | Recover account                     | No          |
| POST   | `/api/users/unblock-account`                             | Request account unblock             | No          |
| PUT    | `/api/users/unblock-account`                             | Reset login error counter           | No          |
| GET    | `/api/users/public-key/{email}`                          | Get public key by email             | No          |
| POST   | `/api/users/attempt-change-email`                        | Initiate email change               | No          |
| POST   | `/api/users/attempt-change-email/{id}/accept`            | Accept email change                 | No          |
| GET    | `/api/users/attempt-change-email/{id}/verify-expiration` | Verify email-change link expiration | No          |
| GET    | `/api/users/meet-token/beta`                             | Get Meet token (beta)               | No          |
| GET    | `/api/users/meet-token/anon`                             | Get anonymous Meet token            | No          |
| POST   | `/api/users/notification-token`                          | Add a notification token            | No          |
| POST   | `/api/users/email-verification/send`                     | Send account verification email     | No          |
| POST   | `/api/users/email-verification`                          | Verify user email                   | No          |
| PATCH  | `/api/users/profile`                                     | Update user profile                 | No          |
| PUT    | `/api/users/avatar`                                      | Upload user avatar                  | No          |
| DELETE | `/api/users/avatar`                                      | Delete user avatar                  | No          |
| POST   | `/api/users/deactivation/send`                           | Send deactivation email             | No          |
| POST   | `/api/users/deactivation/confirm`                        | Confirm user deactivation           | No          |
| GET    | `/api/users/usage`                                       | Get used storage space              | No          |
| GET    | `/api/users/limit`                                       | Get maximum allowed storage         | No          |

### Fuzzy

| Method | Endpoint              | Description                  | Implemented |
| ------ | --------------------- | ---------------------------- | ----------- |
| GET    | `/api/fuzzy/{search}` | Search items by partial name | No          |

### Backup

| Method | Endpoint                              | Description                             | Implemented |
| ------ | ------------------------------------- | --------------------------------------- | ----------- |
| POST   | `/api/backup/activate`                | Activate user backup                    | No          |
| POST   | `/api/backup/deviceAsFolder`          | Create device as folder                 | No          |
| GET    | `/api/backup/deviceAsFolder`          | List device-as-folders                  | No          |
| GET    | `/api/backup/deviceAsFolder/{uuid}`   | Get device-as-folder by UUID            | No          |
| PATCH  | `/api/backup/deviceAsFolder/{uuid}`   | Update device-as-folder                 | No          |
| GET    | `/api/backup/deviceAsFolderById/{id}` | Get device-as-folder by ID (deprecated) | No          |
| GET    | `/api/backup/devices`                 | List user devices                       | No          |
| DELETE | `/api/backup/devices/{deviceId}`      | Delete a user device                    | No          |
| GET    | `/api/backup/{mac}`                   | Get backups for a device                | No          |
| DELETE | `/api/backup/{backupId}`              | Delete a backup                         | No          |

### Trash

| Method | Endpoint                               | Description                              | Implemented |
| ------ | -------------------------------------- | ---------------------------------------- | ----------- |
| GET    | `/api/storage/trash/paginated`         | Get paginated trash content              | No          |
| POST   | `/api/storage/trash/add`               | Add files/folders to trash               | No          |
| DELETE | `/api/storage/trash/all`               | Delete all items from user’s trash       | No          |
| DELETE | `/api/storage/trash/all/request`       | (alternate) Delete all trash via request | No          |
| DELETE | `/api/storage/trash`                   | Delete specified trash items             | No          |
| DELETE | `/api/storage/trash/file/{fileId}`     | Delete single file from trash            | No          |
| DELETE | `/api/storage/trash/folder/{folderId}` | Delete single folder from trash          | No          |

### Sends

| Method | Endpoint              | Description         | Implemented |
| ------ | --------------------- | ------------------- | ----------- |
| POST   | `/api/links`          | Create a send link  | No          |
| GET    | `/api/links/{linkId}` | Get send link by ID | No          |

### Device

| Method | Endpoint                  | Description           | Implemented |
| ------ | ------------------------- | --------------------- | ----------- |
| POST   | `/api/device/context`     | Get device context    | No          |
| POST   | `/api/device/geolocation` | Get geolocation by IP | No          |

### Gateway

| Method | Endpoint                                                      | Description                         | Implemented |
| ------ | ------------------------------------------------------------- | ----------------------------------- | ----------- |
| POST   | `/api/gateway/workspaces`                                     | Initiate a workspace                | No          |
| DELETE | `/api/gateway/workspaces`                                     | Destroy a workspace                 | No          |
| PUT    | `/api/gateway/workspaces/storage`                             | Update a workspace                  | No          |
| POST   | `/api/gateway/workspaces/{workspaceId}/storage/upgrade-check` | Pre-check workspace storage upgrade | No          |
| GET    | `/api/gateway/users`                                          | Get gateway user details            | No          |
| GET    | `/api/gateway/users/storage/stackability`                     | Check if user can expand storage    | No          |
| PATCH  | `/api/gateway/users/{uuid}`                                   | Update user plan and storage        | No          |
