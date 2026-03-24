export interface CreateTourPayload {
    name: string;
    description: string;
    difficulty: string;
    tags: string[];
    price: number;
}
  
export interface UpdateTourPayload {
    name: string;
    description: string;
    difficulty: string;
    tags: string[];
    price: number;
}
