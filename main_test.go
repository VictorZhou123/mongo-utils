package main

import "testing"

func Test_toOBSPath(t *testing.T) {
	type args struct {
		link string
	}
	tests := []struct {
		name        string
		args        args
		wantObspath string
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{
				link: "https://big-model-deploy.obs.cn-central-221.ovaijisuan.com/wukong-huahua/AI-gallery/gallery/上海陆家嘴 未来城市 科幻风格-00.png",
			},
			wantObspath: "wukong-huahua/AI-gallery/gallery/MindSpore/1677053755/上海陆家嘴 未来城市 科幻风格-00.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotObspath := toOBSPath(tt.args.link); gotObspath != tt.wantObspath {
				t.Errorf("toOBSPath() = %v, want %v", gotObspath, tt.wantObspath)
			}
		})
	}
}
