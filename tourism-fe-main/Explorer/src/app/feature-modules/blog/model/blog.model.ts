// Definiše objekat koji šaljemo na backend prilikom kreiranja bloga
export interface CreateBlogPayload {
    title: string;
    content: string; // Ovo je polje za Markdown
    images?: string[]; // Niz Base64 stringova slika, opciono
   // createdAt: string;
  }
  
  // Definiše kompletan Blog objekat koji dobijamo kao odgovor od servera
  export interface Blog {
    id: string;
    title: string;
    content: string;
    htmlContent: string; 
    authorId: number;
    authorUsername? : string;
    createdAt: string;
    updatedAt: string;
    images?: string[];
    comments: BlogComment[]; // Za sada može any, kasnije ćete definisati Comment model
    likes: number[];
  }
  
  export interface BlogComment {
  id: string;
  authorId: number;
  text: string;
  createdAt: string;
  updatedAt: string;
  authorUsername?: string;
}

export interface AddCommentPayload {
  text: string;
}

export interface UpdateBlogPayload {
    title: string;
    content: string; 
    images?: string[]; 
}

export interface UpdateCommentPayload {
    text: string;
}
