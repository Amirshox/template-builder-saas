'use client';

import { useQuery } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeft, Plus } from 'lucide-react';
import { getTemplate, getVersions } from '@/lib/api';

interface TemplateVersion {
    id: string;
    version: number;
    status: string;
    created_at: string;
    template_json: any;
}

export default function TemplateSettingsPage({ params }: { params: { id: string } }) {
    const router = useRouter();
    const templateId = params.id;

    // Fetch Template Details
    const { data: template, isLoading: isTemplateLoading } = useQuery({
        queryKey: ['template', templateId],
        queryFn: () => getTemplate(templateId),
    });

    // Fetch Versions
    const { data: versions, isLoading: isVersionsLoading } = useQuery({
        queryKey: ['versions', templateId],
        queryFn: () => getVersions(templateId),
    });

    if (isTemplateLoading || isVersionsLoading) {
        return <div className="p-8 text-center text-gray-500">Loading...</div>;
    }

    return (
        <main className="min-h-screen bg-gray-50 p-8">
            <div className="mx-auto max-w-4xl">
                <div className="mb-8">
                    <Link href="/" className="flex items-center text-sm text-gray-500 hover:text-gray-900 mb-4">
                        <ArrowLeft className="w-4 h-4 mr-1" />
                        Back to Dashboard
                    </Link>
                    <div className="flex justify-between items-center">
                        <h1 className="text-3xl font-bold text-gray-900">Template Settings</h1>
                    </div>
                </div>

                <div className="bg-white shadow rounded-lg p-6 mb-8">
                    <h2 className="text-xl font-semibold mb-4">Details</h2>

                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Name</label>
                            <div className="mt-1 text-sm text-gray-900">{template?.name}</div>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Type</label>
                            <div className="mt-1 text-sm text-gray-900 capitalize">{template?.type}</div>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Status</label>
                            <div className="mt-1 text-sm text-gray-900 capitalize">{template?.status}</div>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Created At</label>
                            <div className="mt-1 text-sm text-gray-900">{template ? new Date(template.created_at).toLocaleString() : '-'}</div>
                        </div>
                    </div>
                </div>

                <div className="bg-white shadow rounded-lg p-6">
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="text-xl font-semibold">Version History</h2>
                        <Link
                            href={`/editor/${templateId}`}
                            className="flex items-center gap-1 text-sm text-blue-600 hover:text-blue-500 font-medium"
                        >
                            <Plus className="w-4 h-4" /> New Version
                        </Link>
                    </div>
                    <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
                        <table className="min-w-full divide-y divide-gray-300">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th scope="col" className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6">Version</th>
                                    <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Status</th>
                                    <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Created At</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200 bg-white">
                                {versions?.map((version: TemplateVersion) => (
                                    <tr key={version.id}>
                                        <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">
                                            v{version.version}
                                        </td>
                                        <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500 capitalize">{version.status}</td>
                                        <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{new Date(version.created_at).toLocaleDateString()}</td>
                                    </tr>
                                ))}
                                {versions?.length === 0 && (
                                    <tr>
                                        <td colSpan={3} className="text-center py-4 text-gray-500 text-sm">No versions found</td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </main>
    );
}
