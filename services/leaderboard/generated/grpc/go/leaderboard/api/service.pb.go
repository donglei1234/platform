// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: leaderboard/api/service.proto

package pb

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

var File_leaderboard_api_service_proto protoreflect.FileDescriptor

var file_leaderboard_api_service_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6c,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xb9, 0x08, 0x0a, 0x12,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x84, 0x01, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x4b, 0x12, 0x1e,
	0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e,
	0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x4b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26,
	0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e,
	0x47, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x31, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2b, 0x12, 0x29,
	0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f,
	0x7b, 0x61, 0x70, 0x70, 0x49, 0x64, 0x7d, 0x2f, 0x7b, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x7d, 0x12, 0x84, 0x01, 0x0a, 0x07, 0x47, 0x65,
	0x74, 0x4d, 0x54, 0x6f, 0x4e, 0x12, 0x1e, 0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f,
	0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x54, 0x6f, 0x4e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f,
	0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x31, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x2b, 0x12, 0x29, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x7b, 0x61, 0x70, 0x70, 0x49, 0x64, 0x7d, 0x2f, 0x7b,
	0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x7d,
	0x12, 0x88, 0x01, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x49, 0x64, 0x52, 0x61, 0x6e, 0x6b, 0x12, 0x20,
	0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e,
	0x47, 0x65, 0x74, 0x49, 0x64, 0x52, 0x61, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x21, 0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70,
	0x62, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x64, 0x52, 0x61, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x36, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x30, 0x12, 0x2e, 0x2f, 0x76, 0x31,
	0x2f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x7b, 0x61, 0x70,
	0x70, 0x49, 0x64, 0x7d, 0x2f, 0x7b, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x4e, 0x61, 0x6d, 0x65, 0x7d, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x9c, 0x01, 0x0a, 0x12,
	0x47, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x53, 0x69,
	0x7a, 0x65, 0x12, 0x22, 0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x53, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x36, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x30, 0x12, 0x2e, 0x2f, 0x76, 0x31, 0x2f,
	0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x7b, 0x61, 0x70, 0x70,
	0x49, 0x64, 0x7d, 0x2f, 0x7b, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x4e, 0x61, 0x6d, 0x65, 0x7d, 0x2f, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x6e, 0x0a, 0x0b, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x22, 0x2e, 0x6c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e, 0x4e,
	0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1a,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x3a, 0x01, 0x2a, 0x12, 0x89, 0x01, 0x0a, 0x0c, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x20, 0x2e, 0x6c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74,
	0x49, 0x64, 0x52, 0x61, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e, 0x4e,
	0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x36,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x30, 0x2a, 0x2e, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x7b, 0x61, 0x70, 0x70, 0x49, 0x64, 0x7d, 0x2f,
	0x7b, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x4e, 0x61, 0x6d, 0x65,
	0x7d, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x74, 0x0a, 0x0e, 0x4e, 0x65, 0x77, 0x4c, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x12, 0x25, 0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e, 0x4e, 0x65, 0x77, 0x4c, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1f, 0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x62,
	0x2e, 0x4e, 0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x1a, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x3a, 0x01, 0x2a, 0x12, 0x79, 0x0a, 0x10,
	0x52, 0x65, 0x73, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x12, 0x22, 0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70,
	0x62, 0x2e, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x2e, 0x70, 0x62, 0x2e, 0x4e, 0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x22, 0x15, 0x2f,
	0x76, 0x31, 0x2f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x72,
	0x65, 0x73, 0x65, 0x74, 0x3a, 0x01, 0x2a, 0x42, 0x14, 0x5a, 0x12, 0x6c, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_leaderboard_api_service_proto_goTypes = []interface{}{
	(*GetTopKRequest)(nil),             // 0: leaderboard.pb.GetTopKRequest
	(*GetMToNRequest)(nil),             // 1: leaderboard.pb.GetMToNRequest
	(*GetIdRankRequest)(nil),           // 2: leaderboard.pb.GetIdRankRequest
	(*LeaderboardRequest)(nil),         // 3: leaderboard.pb.LeaderboardRequest
	(*UpdateScoreRequest)(nil),         // 4: leaderboard.pb.UpdateScoreRequest
	(*NewLeaderboardRequest)(nil),      // 5: leaderboard.pb.NewLeaderboardRequest
	(*GetLeaderboardResponse)(nil),     // 6: leaderboard.pb.GetLeaderboardResponse
	(*GetIdRankResponse)(nil),          // 7: leaderboard.pb.GetIdRankResponse
	(*GetLeaderBoardSizeResponse)(nil), // 8: leaderboard.pb.GetLeaderBoardSizeResponse
	(*NothingResponse)(nil),            // 9: leaderboard.pb.NothingResponse
}
var file_leaderboard_api_service_proto_depIdxs = []int32{
	0, // 0: leaderboard.pb.LeaderboardService.GetTopK:input_type -> leaderboard.pb.GetTopKRequest
	1, // 1: leaderboard.pb.LeaderboardService.GetMToN:input_type -> leaderboard.pb.GetMToNRequest
	2, // 2: leaderboard.pb.LeaderboardService.GetIdRank:input_type -> leaderboard.pb.GetIdRankRequest
	3, // 3: leaderboard.pb.LeaderboardService.GetLeaderBoardSize:input_type -> leaderboard.pb.LeaderboardRequest
	4, // 4: leaderboard.pb.LeaderboardService.UpdateScore:input_type -> leaderboard.pb.UpdateScoreRequest
	2, // 5: leaderboard.pb.LeaderboardService.DeleteMember:input_type -> leaderboard.pb.GetIdRankRequest
	5, // 6: leaderboard.pb.LeaderboardService.NewLeaderboard:input_type -> leaderboard.pb.NewLeaderboardRequest
	3, // 7: leaderboard.pb.LeaderboardService.ResetLeaderboard:input_type -> leaderboard.pb.LeaderboardRequest
	6, // 8: leaderboard.pb.LeaderboardService.GetTopK:output_type -> leaderboard.pb.GetLeaderboardResponse
	6, // 9: leaderboard.pb.LeaderboardService.GetMToN:output_type -> leaderboard.pb.GetLeaderboardResponse
	7, // 10: leaderboard.pb.LeaderboardService.GetIdRank:output_type -> leaderboard.pb.GetIdRankResponse
	8, // 11: leaderboard.pb.LeaderboardService.GetLeaderBoardSize:output_type -> leaderboard.pb.GetLeaderBoardSizeResponse
	9, // 12: leaderboard.pb.LeaderboardService.UpdateScore:output_type -> leaderboard.pb.NothingResponse
	9, // 13: leaderboard.pb.LeaderboardService.DeleteMember:output_type -> leaderboard.pb.NothingResponse
	9, // 14: leaderboard.pb.LeaderboardService.NewLeaderboard:output_type -> leaderboard.pb.NothingResponse
	9, // 15: leaderboard.pb.LeaderboardService.ResetLeaderboard:output_type -> leaderboard.pb.NothingResponse
	8, // [8:16] is the sub-list for method output_type
	0, // [0:8] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_leaderboard_api_service_proto_init() }
func file_leaderboard_api_service_proto_init() {
	if File_leaderboard_api_service_proto != nil {
		return
	}
	file_leaderboard_api_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_leaderboard_api_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_leaderboard_api_service_proto_goTypes,
		DependencyIndexes: file_leaderboard_api_service_proto_depIdxs,
	}.Build()
	File_leaderboard_api_service_proto = out.File
	file_leaderboard_api_service_proto_rawDesc = nil
	file_leaderboard_api_service_proto_goTypes = nil
	file_leaderboard_api_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LeaderboardServiceClient is the client API for LeaderboardService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LeaderboardServiceClient interface {
	GetTopK(ctx context.Context, in *GetTopKRequest, opts ...grpc.CallOption) (*GetLeaderboardResponse, error)
	GetMToN(ctx context.Context, in *GetMToNRequest, opts ...grpc.CallOption) (*GetLeaderboardResponse, error)
	GetIdRank(ctx context.Context, in *GetIdRankRequest, opts ...grpc.CallOption) (*GetIdRankResponse, error)
	GetLeaderBoardSize(ctx context.Context, in *LeaderboardRequest, opts ...grpc.CallOption) (*GetLeaderBoardSizeResponse, error)
	UpdateScore(ctx context.Context, in *UpdateScoreRequest, opts ...grpc.CallOption) (*NothingResponse, error)
	DeleteMember(ctx context.Context, in *GetIdRankRequest, opts ...grpc.CallOption) (*NothingResponse, error)
	NewLeaderboard(ctx context.Context, in *NewLeaderboardRequest, opts ...grpc.CallOption) (*NothingResponse, error)
	ResetLeaderboard(ctx context.Context, in *LeaderboardRequest, opts ...grpc.CallOption) (*NothingResponse, error)
}

type leaderboardServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLeaderboardServiceClient(cc grpc.ClientConnInterface) LeaderboardServiceClient {
	return &leaderboardServiceClient{cc}
}

func (c *leaderboardServiceClient) GetTopK(ctx context.Context, in *GetTopKRequest, opts ...grpc.CallOption) (*GetLeaderboardResponse, error) {
	out := new(GetLeaderboardResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.pb.LeaderboardService/GetTopK", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardServiceClient) GetMToN(ctx context.Context, in *GetMToNRequest, opts ...grpc.CallOption) (*GetLeaderboardResponse, error) {
	out := new(GetLeaderboardResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.pb.LeaderboardService/GetMToN", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardServiceClient) GetIdRank(ctx context.Context, in *GetIdRankRequest, opts ...grpc.CallOption) (*GetIdRankResponse, error) {
	out := new(GetIdRankResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.pb.LeaderboardService/GetIdRank", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardServiceClient) GetLeaderBoardSize(ctx context.Context, in *LeaderboardRequest, opts ...grpc.CallOption) (*GetLeaderBoardSizeResponse, error) {
	out := new(GetLeaderBoardSizeResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.pb.LeaderboardService/GetLeaderBoardSize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardServiceClient) UpdateScore(ctx context.Context, in *UpdateScoreRequest, opts ...grpc.CallOption) (*NothingResponse, error) {
	out := new(NothingResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.pb.LeaderboardService/UpdateScore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardServiceClient) DeleteMember(ctx context.Context, in *GetIdRankRequest, opts ...grpc.CallOption) (*NothingResponse, error) {
	out := new(NothingResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.pb.LeaderboardService/DeleteMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardServiceClient) NewLeaderboard(ctx context.Context, in *NewLeaderboardRequest, opts ...grpc.CallOption) (*NothingResponse, error) {
	out := new(NothingResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.pb.LeaderboardService/NewLeaderboard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardServiceClient) ResetLeaderboard(ctx context.Context, in *LeaderboardRequest, opts ...grpc.CallOption) (*NothingResponse, error) {
	out := new(NothingResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.pb.LeaderboardService/ResetLeaderboard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LeaderboardServiceServer is the server API for LeaderboardService service.
type LeaderboardServiceServer interface {
	GetTopK(context.Context, *GetTopKRequest) (*GetLeaderboardResponse, error)
	GetMToN(context.Context, *GetMToNRequest) (*GetLeaderboardResponse, error)
	GetIdRank(context.Context, *GetIdRankRequest) (*GetIdRankResponse, error)
	GetLeaderBoardSize(context.Context, *LeaderboardRequest) (*GetLeaderBoardSizeResponse, error)
	UpdateScore(context.Context, *UpdateScoreRequest) (*NothingResponse, error)
	DeleteMember(context.Context, *GetIdRankRequest) (*NothingResponse, error)
	NewLeaderboard(context.Context, *NewLeaderboardRequest) (*NothingResponse, error)
	ResetLeaderboard(context.Context, *LeaderboardRequest) (*NothingResponse, error)
}

// UnimplementedLeaderboardServiceServer can be embedded to have forward compatible implementations.
type UnimplementedLeaderboardServiceServer struct {
}

func (*UnimplementedLeaderboardServiceServer) GetTopK(context.Context, *GetTopKRequest) (*GetLeaderboardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopK not implemented")
}
func (*UnimplementedLeaderboardServiceServer) GetMToN(context.Context, *GetMToNRequest) (*GetLeaderboardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMToN not implemented")
}
func (*UnimplementedLeaderboardServiceServer) GetIdRank(context.Context, *GetIdRankRequest) (*GetIdRankResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIdRank not implemented")
}
func (*UnimplementedLeaderboardServiceServer) GetLeaderBoardSize(context.Context, *LeaderboardRequest) (*GetLeaderBoardSizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeaderBoardSize not implemented")
}
func (*UnimplementedLeaderboardServiceServer) UpdateScore(context.Context, *UpdateScoreRequest) (*NothingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateScore not implemented")
}
func (*UnimplementedLeaderboardServiceServer) DeleteMember(context.Context, *GetIdRankRequest) (*NothingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMember not implemented")
}
func (*UnimplementedLeaderboardServiceServer) NewLeaderboard(context.Context, *NewLeaderboardRequest) (*NothingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewLeaderboard not implemented")
}
func (*UnimplementedLeaderboardServiceServer) ResetLeaderboard(context.Context, *LeaderboardRequest) (*NothingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetLeaderboard not implemented")
}

func RegisterLeaderboardServiceServer(s *grpc.Server, srv LeaderboardServiceServer) {
	s.RegisterService(&_LeaderboardService_serviceDesc, srv)
}

func _LeaderboardService_GetTopK_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTopKRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).GetTopK(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.pb.LeaderboardService/GetTopK",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).GetTopK(ctx, req.(*GetTopKRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderboardService_GetMToN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMToNRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).GetMToN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.pb.LeaderboardService/GetMToN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).GetMToN(ctx, req.(*GetMToNRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderboardService_GetIdRank_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIdRankRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).GetIdRank(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.pb.LeaderboardService/GetIdRank",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).GetIdRank(ctx, req.(*GetIdRankRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderboardService_GetLeaderBoardSize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaderboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).GetLeaderBoardSize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.pb.LeaderboardService/GetLeaderBoardSize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).GetLeaderBoardSize(ctx, req.(*LeaderboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderboardService_UpdateScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).UpdateScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.pb.LeaderboardService/UpdateScore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).UpdateScore(ctx, req.(*UpdateScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderboardService_DeleteMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIdRankRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).DeleteMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.pb.LeaderboardService/DeleteMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).DeleteMember(ctx, req.(*GetIdRankRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderboardService_NewLeaderboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewLeaderboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).NewLeaderboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.pb.LeaderboardService/NewLeaderboard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).NewLeaderboard(ctx, req.(*NewLeaderboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderboardService_ResetLeaderboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaderboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).ResetLeaderboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.pb.LeaderboardService/ResetLeaderboard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).ResetLeaderboard(ctx, req.(*LeaderboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LeaderboardService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "leaderboard.pb.LeaderboardService",
	HandlerType: (*LeaderboardServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTopK",
			Handler:    _LeaderboardService_GetTopK_Handler,
		},
		{
			MethodName: "GetMToN",
			Handler:    _LeaderboardService_GetMToN_Handler,
		},
		{
			MethodName: "GetIdRank",
			Handler:    _LeaderboardService_GetIdRank_Handler,
		},
		{
			MethodName: "GetLeaderBoardSize",
			Handler:    _LeaderboardService_GetLeaderBoardSize_Handler,
		},
		{
			MethodName: "UpdateScore",
			Handler:    _LeaderboardService_UpdateScore_Handler,
		},
		{
			MethodName: "DeleteMember",
			Handler:    _LeaderboardService_DeleteMember_Handler,
		},
		{
			MethodName: "NewLeaderboard",
			Handler:    _LeaderboardService_NewLeaderboard_Handler,
		},
		{
			MethodName: "ResetLeaderboard",
			Handler:    _LeaderboardService_ResetLeaderboard_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "leaderboard/api/service.proto",
}
