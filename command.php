#!/usr/bin/env php
<?php
// This contains first argument
$message = $argv[1];

// Decode to get original value
$original = base64_decode($message);

// Start processing
print_r(date("Y-m-d H:i:s"));

file_put_contents("ret.log", $original."\n", FILE_APPEND);
echo "\n";
if (true) {
    print_r($original);

    // All well, then return 0
    exit(0);
}

// Let rabbitmq-cli-consumer know someting went wrong, message will be requeued.
exit(1);