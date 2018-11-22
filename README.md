[![Build Status](https://travis-ci.com/nupplaphil/kopano-ldap.svg?branch=master)](https://travis-ci.com/nupplaphil/kopano-ldap)
[![codecov](https://codecov.io/gh/nupplaphil/kopano-ldap/branch/master/graph/badge.svg)](https://codecov.io/gh/nupplaphil/kopano-ldap)
[![Go Report Card](https://goreportcard.com/badge/github.com/nupplaphil/kopano-ldap)](https://goreportcard.com/report/github.com/nupplaphil/kopano-ldap)

# LDAP Command Line Interface for Kopano

This CLI is built for the User Management of [Kopano](https://kopano.io).
It is inspired by `kopano-cli` and is written in [Go](https://golang.org).

It should be usef in combination with the general `kopano-cli` tool. `kopano-cli` starts, where `kopano-cli` stops with LDAP support.

## General usage of `kopano-ld` tool

The `kopano-ld` is an administration tool for managing user and groups in LDAP. 
If you use an other backend than `ldap` for Kopano (f.e `DB` or `unix`), please look for the [8. User Management](https://documentation.kopano.io/kopanocore_administrator_manual/user_management.html) in Kopano.

The tool can be used to get more information about users and groups too.

### Listing users

All available users or groups can be displayed by using the following commands:

```bash
> kopano-ld user --list
> kopano-ld group --list
```

## Users

### Displaying details

To display more information of a specific user, use:
```bash
> kopano-ld user john
Name:                     johndoe
Full name:                John Doe
Email address:            john@doe.com
Active:                   yes
Administrator:            no
Features Enabled:         mobile
Features Disabled:        imap; pop3
```

To display more information of a specific group, use:

```bash
TO BE DEFINED (TODO)
```

### Creating users with LDAP

To create a new user, use the following command:

```bash
> kopano-ld user create --user <user name> \
                      --password <password> \
                      --email <email> \
                      --fullname <full name> \
                      --active <active> \
                      --admin-evel <administrator>
```

The fields between <> should be filled in as follows:

   - `User name`: The name of the user. With this name the user will log on to the store.
   - `Password`: The password in plain text. The password will be stored encrypted in LDAP.
   - `Email`: The email address of the user. Often this is `<user name>@<email domain>`. 
              You can define more than one email address, which will be set as an alias for this user
   - `Active`: The active state of this user. If set to `yes` (or not set), the user is able to login, otherwise `no`
   - `Administrator`: This value should be `0` or `1`. When a use is administrator, the user will be allowed to open all Kopano stores of any user.
                      **(TODO - not working yet)**
                      
### Updating user information with LDAP

The same `kopano-ld` tool can be used to update LDAP information. Use the following command to update:

```bash
TO BE DEFINED (TODO)
```

### Deleting users with LDAP

To delete a user from LDAP, use the following command:

```bash
> kopano-ld user delete --user <user name>
```

The user will be deleted from LDAP. However the store will be kept in the database.

## Groups

TO BE DEFINED (TODO)

## Creating Groups with LDAP

TO BE DEFINED (TODO)

## Deleting Groups with LDAP

TO BE DEFINED (TODO)

## Kopano Feature management

Some features within KC can be disabled. 
By default, all features are disabled. 
Enabling can be done globally or on a per-user basis. 
When a feature has been globally disabled, you may enable the feature in a per-user basis too. 

Currently the only features that can be controlled are **`imap`**, **`pop3`** and **`mobile`**.

If the **`pop3`** feature is disabled, users won’t be able to login using the POP3 protocol.
The same goes for the **`imap`** feature, but this has an extra effect aswell. 
When a user receives email when the **`imap`** feature is enabled, the original email and some other imap optimized data will also be saved in the Kopano database and attachment directory.
This will make the IMAP services provided by the kopano-gateway more reliable. 
On the other hand, it will also use more diskspace. 

Disabling the **`imap`** feature will thus save diskspace.

### Globally enabling features

see [8.7.1. Globally enabling features](https://documentation.kopano.io/kopanocore_administrator_manual/user_management.html#globally-enabling-features)

### Per-user en- or disabling features

Managing the feature per user in LDAP, the `kopano-ld` tool has to be used to control the features:

```bash
> kopano-ld user feature add --user john --add imap
> kopano-ld user feature rem --user john --rem mobile
``` 

In LDAP, the features will be managed from the two attributes `kopanoEnabledFeatures` and `kopanoDisabledFeatures`.
