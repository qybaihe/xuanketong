// 类型声明文件 for mock-upload-api.js

export interface MaterialData {
  materialName: string;
  type: 'cloud' | 'note';
  service: string;
  link: string;
  accessPassword?: string;
  description?: string;
}

export interface CourseMaterial {
  id: number;
  courseId: number;
  materialName: string;
  type: 'cloud' | 'note';
  service: string;
  link: string;
  accessPassword?: string;
  description?: string;
  uploaderId: number;
  uploaderName: string;
  uploadTime: string;
  downloadCount: number;
}

export interface Note {
  id: number;
  courseId: number;
  content: string;
  authorId: number;
  authorName: string;
  createdAt: string;
  isOwner: boolean;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
}

export interface MaterialsResponse {
  total: number;
  materials: CourseMaterial[];
}

export interface NotesResponse {
  notes: Note[];
}

export interface SaveNoteResponse {
  id: number;
  courseId: number;
  content: string;
  lastModified: string;
}

export declare function uploadMaterial(courseId: number, materialData: MaterialData): Promise<ApiResponse<CourseMaterial>>;

export declare function getCourseMaterials(courseId: number): Promise<ApiResponse<MaterialsResponse>>;

export declare function saveCourseNote(courseId: number, content: string): Promise<ApiResponse<SaveNoteResponse>>;

export declare function getCourseNotes(courseId: number): Promise<ApiResponse<NotesResponse>>;
