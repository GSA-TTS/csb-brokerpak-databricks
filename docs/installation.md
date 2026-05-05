# Installing the broker for Databricks

The broker service and the Databricks brokerpak can be pushed and registered on a Cloud Foundry foundation.

Documentation for broker configuration can be found [here](./configuration.md).

## Requirements

### Databricks Account

A Databricks workspace is required to provision services. The brokerpak needs:

- `DATABRICKS_HOST` - the workspace URL (e.g. `https://adb-123456789.azuredatabricks.net`)
- `DATABRICKS_TOKEN` - a personal access token with permission to manage clusters

#### Creating a Personal Access Token

1. Log in to your Databricks workspace.
2. Click your user name in the top right corner and select **User Settings**.
3. Go to **Access tokens** and click **Generate new token**.
4. Enter a comment and click **Generate**.
5. Save the token value securely.

### MySQL Database for Broker State

The broker keeps service instance and binding information in a MySQL database.

#### Binding a MySQL Database

If there is an existing broker in the foundation that can provision a MySQL instance use `cf create-service`
to create a new MySQL instance. Then use `cf bind-service` to bind that instance to the service broker.

## Step By Step From a Pre-built Release

### Fetch A Broker and Databricks Brokerpak

Download a Cloud Service Broker release from https://github.com/cloudfoundry/cloud-service-broker/releases.
Find the latest release matching the name pattern `vX.X.X`.
Change filename `cloud-service-broker.linux` to `cloud-service-broker`.
Add execution permissions `chmod +x cloud-service-broker`

Download a Databricks Brokerpak release from this repository's releases.
Find the latest release matching the name pattern `X.X.X`.

Put the `cloud-service-broker` and `databricks-services-X.X.X.brokerpak` into the same directory on your workstation.

### Build Config File

To avoid putting any sensitive information in environment variables, a config file can be used.

Create a file named `config.yml` in the same directory the broker and brokerpak have been downloaded to. Its contents should be:

```yaml
databricks:
  host: https://adb-123456789.azuredatabricks.net
  token: your-databricks-personal-access-token
```

### Push and Register the Broker

Push the broker as a binary application:

```bash
SECURITY_USER_NAME=someusername
SECURITY_USER_PASSWORD=somepassword
APP_NAME=cloud-service-broker-databricks

chmod +x cloud-service-broker
cf push "${APP_NAME}" -c './cloud-service-broker serve --config config.yml' -b binary_buildpack --random-route --no-start
```

Bind the MySQL database and start the service broker:

```bash
cf bind-service cloud-service-broker-databricks csb-sql
cf start "${APP_NAME}"
```

Register the service broker:

```bash
BROKER_NAME=csb-$USER

cf create-service-broker "${BROKER_NAME}" "${SECURITY_USER_NAME}" "${SECURITY_USER_PASSWORD}" https://$(LANG=EN cf app "${APP_NAME}" | grep 'routes:' | cut -d ':' -f 2 | xargs) --space-scoped || cf update-service-broker "${BROKER_NAME}" "${SECURITY_USER_NAME}" "${SECURITY_USER_PASSWORD}" https://$(LANG=EN cf app "${APP_NAME}" | grep 'routes:' | cut -d ':' -f 2 | xargs)
```

Once this completes, the output from `cf marketplace` should include:

```
csb-databricks-workspace    default   Databricks workspace cluster with default configuration
```

## Uninstalling the Broker

First, make sure there are all service instances created with `cf create-service` have been destroyed
with `cf delete-service` otherwise removing the broker will fail.

### Unregister the Broker

```bash
cf delete-service-broker csb-$USER
```

### Uninstall the Broker

```bash
cf delete cloud-service-broker-databricks
```
