'use client';

import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { fetchTemplates, logout } from '@/lib/api';

export default function Home() {
  const router = useRouter();

  // Auth Check
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
    }
  }, [router]);

  const { data: templates, isLoading, error } = useQuery({
    queryKey: ['templates'],
    queryFn: () => fetchTemplates(""), // OrgId is now handled by backend token
    retry: false,
  });

  return (
    <main className="min-h-screen bg-gray-50 p-8">
      <div className="mx-auto max-w-7xl">
        <div className="mb-8 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">My Templates</h1>
            <p className="mt-1 text-sm text-gray-500">Manage and create document templates</p>
          </div>
          <div className="flex gap-4">
            <button
              onClick={logout}
              className="rounded-md bg-white px-4 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
            >
              Sign out
            </button>
            <Link
              href="/editor/new"
              className="rounded-md bg-blue-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
            >
              New Template
            </Link>
          </div>
        </div>

        {error && (
          <div className="mb-8 rounded-md bg-red-50 p-4">
            <div className="flex">
              <div className="ml-3">
                <h3 className="text-sm font-medium text-red-800">Error loading templates</h3>
                <div className="mt-2 text-sm text-red-700">
                  <p>{(error as Error).message || "Unauthorized or Server Error"}</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {isLoading ? (
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-48 animate-pulse rounded-lg bg-gray-200"></div>
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {templates?.map((template) => (
              <div key={template.id} className="block group relative">
                <div className="overflow-hidden rounded-lg bg-white shadow transition-shadow hover:shadow-md">
                  <div className="aspect-video w-full bg-gray-100 object-cover group-hover:bg-gray-200 transition-colors" />
                  <div className="p-4">
                    <h3 className="text-lg font-medium text-gray-900">
                      <Link href={`/editor/${template.id}`} className="after:absolute after:inset-0">
                        {template.name}
                      </Link>
                    </h3>
                    <p className="text-sm text-gray-500 mt-1 capitalize">{template.type} Template</p>
                    <div className="mt-4 flex items-center justify-between text-xs text-gray-400">
                      <span>{new Date(template.created_at).toLocaleDateString()}</span>
                      <span className={`inline-flex items-center rounded-full px-2 py-1 font-medium ${template.status === 'active' ? 'bg-green-50 text-green-700' : 'bg-gray-50 text-gray-600'
                        }`}>
                        {template.status}
                      </span>
                    </div>
                    <div className="mt-2 flex gap-2 relative z-10">
                      <Link
                        href={`/templates/${template.id}/settings`}
                        className="text-xs text-blue-600 hover:text-blue-500"
                      >
                        Settings
                      </Link>
                    </div>
                  </div>
                </div>
              </div>
            ))}

            {templates?.length === 0 && (
              <div className="col-span-full py-12 text-center">
                <p className="text-gray-500">No templates found. Create your first one!</p>
              </div>
            )}
          </div>
        )}
      </div>
    </main>
  );
}
