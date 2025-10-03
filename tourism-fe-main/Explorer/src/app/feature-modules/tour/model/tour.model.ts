export interface Tour {
  id: number;
  authorId: number;
  name: string;
  description: string;
  difficulty: 'Easy' | 'Medium' | 'Hard' | 'Expert';
  tags: string[];
  status: 'Draft' | 'Published' | 'Archived';
  price: number;
  distance?: number;
  publishedAt?: string;
  archivedAt?: string;
  createdAt: string;
  updatedAt: string;
  durations?: TourDuration[];
  keyPoints?: KeyPoint[];
}

export interface KeyPoint {
  id: number;
  tourId: number;
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  address?: string;
  image?: string;
  order: number;
}

export interface TourDuration {
  id: number;
  tourId: number;
  transportType: 'walking' | 'bicycle' | 'car';
  durationMin: number;
}

export interface AddDurationPayload {
  transportType: 'walking' | 'bicycle' | 'car';
  durationMin: number;
}