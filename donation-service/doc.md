sudo docker build -t gcr.io/crowdfunding-460613/donation-service .

docker push gcr.io/crowdfunding-460613/donation-service

gcloud container images list --repository=gcr.io/crowdfunding-460613

gcloud run deploy donation-service \
  --image gcr.io/crowdfunding-460613/donation-service \
  --platform managed \
  --region asia-southeast2 \
  --allow-unauthenticated \
  --port 50051

 https://donation-service-273575294549.asia-southeast2.run.app