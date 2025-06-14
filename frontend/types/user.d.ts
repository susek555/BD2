export interface User {
  accessToken: string;
  refreshToken: string;
  username: string;
  email: string;
  selector: "P" | "C";
  person_name?: string;
  person_surname?: string;
  company_name?: string;
  company_surname?: string;
  errors?: string[];
}
