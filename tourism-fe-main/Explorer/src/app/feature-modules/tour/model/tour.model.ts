export interface Tour {
  id: number;
  authorId: number;
  name: string;
  description: string;
  difficulty: 'Easy' | 'Medium' | 'Hard' | 'Expert';
  tags: string[];
  status: 'Draft' | 'Published' | 'Archived';
  price: number;
  createdAt: string;
  updatedAt: string;
}