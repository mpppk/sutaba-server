#!/usr/bin/env bash
gcloud beta run deploy sutaba-server \
    --image asia.gcr.io/sutaba/sutaba-server \
    --platform managed \
    --region asia-northeast1 \
    --update-env-vars CLASSIFIER_SERVER_HOST=https://classifier-server-lkui2qyzba-an.a.run.app \
    --verbosity debug
