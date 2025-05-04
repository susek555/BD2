import SideBar from "@/app/ui/(home)/sidebar";

export default function Home() {
  return (
  <main>
    <div className="flex flex-grow flex-col md:flex-row"></div>
        <div className="w-full flex-none md:w-64">
          <SideBar />
        </div>
        <div className="flex-grow p-6 md:overflow-y-auto md:p-12">
          <p>Main Page</p>
        </div>
  </main>
  );
}
