# Stalof

Stalof or Stackdriver log flatter is dead simple cli tool that enables easier handling of
gcp stackdriver logs. Stackdriver by default logs in json, and for some cases we only
want/need original logs text entries. For example if you dump mysql slow query log
to sink (google storage bucket), to analyze it with mysqldumpslow you need only textual
part of stackdriver log. You can use stalof to flatten down your log and analyze it.

## Configuration
Stalof uses golang google client library. Auth to gcp is made using default credentials stored in GOOGLE_APPLICATION_CREDENTIALS environment variable.

We suggest creating new service account that will have only read permissions on bucket where you store your logs, and to export this service account .json file path to GOOGLE_APPLICATION_CREDENTIALS env variable.

## Usage

* ```$ stalof your-log-bucket``` - Will read logs from specified bucket and display it back on stdOut. You can pipe this to file or some other tool. Logs will be written to stdOut and potential errors to stdError so be careful when redirecting output.
