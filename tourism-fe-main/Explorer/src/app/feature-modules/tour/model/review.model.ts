export interface Review {
  id: number;
  tourId: number;
  touristId: number;
  touristName?: string;
  touristUsername?: string;
  rating: number; // 1-5
  comment: string;
  visitDate: string;
  images: string[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateReviewRequest {
  tourId: number;
  rating: number;
  comment: string;
  visitDate: string;
  images: string[];
}

export interface UpdateReviewRequest {
  rating?: number;
  comment?: string;
  visitDate?: string;
  images?: string[];
}

export interface ReviewStats {
  averageRating: number;
  reviewCount: number;
}
