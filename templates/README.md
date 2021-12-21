# Templates for the AWS SDK Code Samples

This folder contains template readme files for each of the levels and types of example code found within:

* `sdk.md`: Standard template for language-level content. You probably won't have to use this unless you're Very Special.
* `service.md`: Standard template for a specific service using a specific language. If you're the first one to add examples
  for your service in your language, this is the one for you.
* `sdk-cross-top.md`: Standard template for SDK-level cross-service example content. You probably won't have to use this unless you're
   Very Special.
* `sdk-cross-example.md`: Standard template for a cross-service example. If you're adding a new sample that utilizes multiple services,
   this is the readme to use for that.

## Using these readmes

* When adding a new service's worth of examples, copy the `service.md` file into your service's appropriate folder.
  This should be a part of your PR, and should be the second thing you do other than write the examples.
* When adding a new cross-service example, copy the `sdk-cross-example.md` file into your example's folder as `README.md`.
  This should be a part of your PR, and should be the second thing you do other than write the example.

Each example has a set of placeholders. Fill out the placeholders to the best of your ability. 
