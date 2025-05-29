sudo docker build -t gcr.io/crowdfunding-460613/user-service .

docker push gcr.io/crowdfunding-460613/user-service

gcloud container images list --repository=gcr.io/crowdfunding-460613

gcloud run deploy user-service \
  --image gcr.io/crowdfunding-460613/user-service \
  --platform managed \
  --region asia-southeast2 \
  --allow-unauthenticated \
  --port 50051

 https://user-service-273575294549.asia-southeast2.run.app