/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0.
 */

// snippet-start:[rust.example_code.s3.scenario_getting_started.bin]

// snippet-start:[rust.example_code.s3.basics.imports]

use aws_config::meta::region::RegionProviderChain;
use aws_sdk_s3::{Client, Error, Region};
use uuid::Uuid;

// snippet-end:[rust.example_code.s3.basics.imports]

// snippet-start:[rust.example_code.s3.basics.main]
#[tokio::main]
async fn main() -> Result<(), Error> {
    let (region, client, bucket_name, file_name, key, target_key) = initialize_variables().await;

    if let Err(e) = run_s3_operations(region, client, bucket_name, file_name, key, target_key).await
    {
        println!("{:?}", e);
    };

    Ok(())
}
// snippet-end:[rust.example_code.s3.basics.main]

// snippet-start:[rust.example_code.s3.basics.initialize_variables]
async fn initialize_variables() -> (Region, Client, String, String, String, String) {
    let region_provider = RegionProviderChain::first_try(Region::new("us-west-2"));
    let region = region_provider.region().await.unwrap();

    let shared_config = aws_config::from_env().region(region_provider).load().await;
    let client = Client::new(&shared_config);

    let bucket_name = format!("{}{}", "doc-example-bucket-", Uuid::new_v4().to_string());
    let file_name = "s3/testfile.txt".to_string();
    let key = "test file key name".to_string();
    let target_key = "target_key".to_string();

    (region, client, bucket_name, file_name, key, target_key)
}
// snippet-end:[rust.example_code.s3.basics.initialize_variables]

// snippet-start:[rust.example_code.s3.basics.run_s3_operations]
async fn run_s3_operations(
    region: Region,
    client: Client,
    bucket_name: String,
    file_name: String,
    key: String,
    target_key: String,
) -> Result<(), Error> {
    s3_service::create_bucket(&client, &bucket_name, region.as_ref()).await?;
    s3_service::upload_object(&client, &bucket_name, &file_name, &key).await?;
    s3_service::download_object(&client, &bucket_name, &key).await?;
    s3_service::copy_object(&client, &bucket_name, &key, &target_key).await?;
    s3_service::list_objects(&client, &bucket_name).await?;
    s3_service::delete_objects(&client, &bucket_name).await?;
    s3_service::delete_bucket(&client, &bucket_name).await?;

    Ok(())
}
// snippet-end:[rust.example_code.s3.basics.run_s3_operations]

// snippet-end:[rust.example_code.s3.scenario_getting_started.bin]
