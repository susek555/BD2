import "next-auth";
import { JWT } from "next-auth/jwt";


declare module "next-auth" {
  interface Session {
    user: {
      accessToken: string;
      refreshToken: string;
      username: string;
      email: string;
      selector: "P" | "C";
      personName?: string;
      personSurname?: string;
      companyName?: string;
      companySurname?: string;
      accessTokenExpires: number;
      errors?: string[];
    };
  }
  interface User {
    accessToken: string;
    refreshToken: string;
    username: string;
    email: string;
    selector: "P" | "C";
    personName?: string;
    personSurname?: string;
    companyName?: string;
    companySurname?: string;
    errors?: string[];
  }
}

interface ExtendedJWT extends JWT {
  accessToken: string;
  refreshToken: string;
  username: string;
  email: string;
  selector: "P" | "C";
  personName?: string;
  personSurname?: string;
  companyName?: string;
  companySurname?: string;
  accessTokenExpires: number;
  error?: string;
}
