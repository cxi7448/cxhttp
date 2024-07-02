package faceswap

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/xhttp"
	"strings"
)

type Akool struct {
	Api
	Token        string
	ClientId     string
	ClientSecret string
	UserId       string
	WebhookUrl   string
}

var akool = &Akool{}

func InitAkool(clientId, clientSecret, user_id, webhookUrl string) {
	akool.ClientId = clientId
	akool.ClientSecret = clientSecret
	akool.UserId = user_id
	akool.WebhookUrl = webhookUrl
	clLog.Info("设置clientId[%v],clientSecret[%v],userId[%v]", clientId, clientSecret, user_id)
}

func (this *Akool) GenToken() (string, error) {
	url := "https://openapi.akool.com/api/open/v3/getToken"
	client := xhttp.New(url)
	result := clJson.M{}
	err := client.Post(clJson.M{
		"clientId":     this.ClientId,
		"clientSecret": this.ClientSecret,
	}, &result)
	if err != nil {
		clLog.Error("错误：%v", err)
		return "", err
	}
	if result.Uint32("code") != 1000 || result.Get("token") == "" {
		clLog.Error("错误内容:%v", result.Get("message"))
		return "", fmt.Errorf(result.Get("message"))
	} else {
		this.Token = result.Get("token")
	}
	return this.Token, nil
}

func (this *Akool) FaceSwap(src, face Img) (string, string, error) {
	token, err := this.GenToken()
	if err != nil {
		clLog.Error("生成访问密钥错误:%v", err)
		return "", "", err
	}
	url := "https://openapi.akool.com/api/open/v3/faceswap/highquality/specifyimage"
	client := xhttp.New(url)
	client.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", token),
	})
	result := struct {
		Code uint32 `json:"code"`
		Data struct {
			Id    string `json:"_id"`
			JobId string `json:"job_id"`
			Url   string `json:"url"`
		} `json:"data"`
		Msg string `json:"msg"`
	}{}
	err = client.Post(clJson.M{
		"targetImage": clJson.A{
			clJson.M{
				"path": src.Image,
				"opts": src.Opts,
			},
		},
		"sourceImage": clJson.A{
			clJson.M{
				"path": face.Image,
				"opts": face.Opts,
			},
		},
		"face_enhance": 0,
		"modifyImage":  src.Image,
		"webhookUrl":   this.WebhookUrl,
	}, &result)
	if err != nil {
		clLog.Error("访问[%v]失败:%v", url, err)
		return "", "", err
	}
	clLog.Debug("请求结果:%+v", result)
	if result.Code != 1000 {
		return "", "", fmt.Errorf(result.Msg)
	}
	// 16:27:48 akool.go:65[Err] 请求结果:map[code:1000 data:map[_id:667d2284dca9e468ba8ead23 job_id:20240627082748003-5746 url:https://d2qf6ukcym4kn9.cloudfront.net/final_bdd1c994c4cd7a58926088ae8a479168-1705462506461-1966-3d389dcf-f9f7-4134-9594-9fc2a0fcc6f4-2272.jpeg] msg:Please be patient! If your results are not generated in three hours, please check your input image.]
	return result.Data.Url, result.Data.Id, err
}

func (this *Akool) FaceSwapVideo(srcs, faces []Img, video_url string) (string, string, error) {
	token, err := this.GenToken()
	if err != nil {
		clLog.Error("生成访问密钥错误:%v", err)
		return "", "", err
	}
	url := "https://openapi.akool.com/api/open/v3/faceswap/highquality/specifyvideo"
	client := xhttp.New(url)
	client.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", token),
	})
	result := struct {
		Code uint32 `json:"code"`
		Data struct {
			Id    string `json:"_id"`
			JobId string `json:"job_id"`
			Url   string `json:"url"`
		} `json:"data"`
		Msg string `json:"msg"`
	}{}
	targetImage := clJson.A{}
	sourceImage := clJson.A{}
	for _, val := range srcs {
		targetImage = append(targetImage, clJson.M{
			"path": val.Image,
			"opts": val.Opts,
		})
	}
	for _, val := range faces {
		sourceImage = append(sourceImage, clJson.M{
			"path": val.Image,
			"opts": val.Opts,
		})
	}
	err = client.Post(clJson.M{
		"targetImage":  targetImage,
		"sourceImage":  sourceImage,
		"face_enhance": 0,
		"modifyVideo":  video_url,
		"webhookUrl":   this.WebhookUrl,
	}, &result)
	if err != nil {
		clLog.Error("访问[%v]失败:%v", url, err)
		return "", "", err
	}
	clLog.Debug("视频请求结果:%+v", result)
	if result.Code != 1000 {
		return "", "", fmt.Errorf(result.Msg)
	}
	// 16:27:48 akool.go:65[Err] 请求结果:map[code:1000 data:map[_id:667d2284dca9e468ba8ead23 job_id:20240627082748003-5746 url:https://d2qf6ukcym4kn9.cloudfront.net/final_bdd1c994c4cd7a58926088ae8a479168-1705462506461-1966-3d389dcf-f9f7-4134-9594-9fc2a0fcc6f4-2272.jpeg] msg:Please be patient! If your results are not generated in three hours, please check your input image.]
	return result.Data.Url, result.Data.Id, err
}

// 0等待中  1成功 2失败
func (this *Akool) CheckResult(id string) (uint32, error) {
	// 16:50:22 akool.go:145[Debug] 视频请求结果:{Code:1000 Data:{Id:6683bf4ddca9e468ba90e360 JobId:20240702085021960-9426 Url:https://d2qf6ukcym4kn9.cloudfront.net/final_26fd899ff32d6c40f880c67f5e1f4519053e47e66b076531a17c59-627da158-1262-4e38-8162-a0e2f4f0b1d9-2554.mp4} Msg:Please be patient! If your results are not generated in three hours, please check your input video.}
	token, err := this.GenToken()
	if err != nil {
		clLog.Error("生成访问密钥错误:%v", err)
		return 0, err
	}
	url := fmt.Sprintf("https://openapi.akool.com/api/open/v3/faceswap/result/listbyids?_ids=%v", id)
	client := xhttp.New(url)
	client.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", token),
	})
	result := struct {
		Code uint32 `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Result struct {
				FaceswapStatus uint32 `json:"faceswap_status"`
				Id             string `json:"_id"`
				Url            string `json:"url"`
				Image          uint32 `json:"image"`
				CreatedAt      string `json:"createdAt"`
				JobId          string `json:"job_id"`
			} `json:"result"`
		} `json:"data"`
	}{}
	err = client.Get(&result)
	if err != nil {
		clLog.Error("访问API失败:%v", err)
		return 0, err
	}
	clLog.Info("请求结果:%+v", result)
	// 16:33:08 akool.go:82[Info] 请求结果:map[code:1000 data:map[result:[map[_id:667d2284dca9e468ba8ead23 createdAt:2024-06-27T08:27:48.006Z deduction_duration:0 faceswap_status:3 image:1 job_id:20240627082748003-5746 uid:2.839571e+06 url:https://d2qf6ukcym4kn9.cloudfront.net/final_bdd1c994c4cd7a58926088ae8a479168-1705462506461-1966-3d389dcf-f9f7-4134-9594-9fc2a0fcc6f4-2272.jpeg video_duration:0]]] msg:OK]
	if result.Code != 1000 {
		clLog.Error("错误:%v", result.Msg)
		return 0, fmt.Errorf("错误内容:%v", result.Msg)
	}
	status := result.Data.Result.FaceswapStatus
	if status < 3 {
		return 0, nil
	}
	if status == 3 {
		return 1, nil
	} else {
		clLog.Error("制作失败:%+v", result)
		return 2, nil
	}
}

func (this *Akool) Detect(image string) (string, error) {
	token, err := this.GenToken()
	if err != nil {
		return "", err
	}
	url := "https://sg3.akool.com/detect"
	client := xhttp.New(url)
	client.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", token),
	})
	//folder := "./tmp/ai/detect"
	//os.MkdirAll(folder, 0700)
	//// 下载图片，并生成base64
	//local_path := fmt.Sprintf("%v/%v", folder, clCommon.Md5([]byte(image)))
	//err = clFile.Download(image, local_path)
	//if err != nil {
	//	clLog.Error("文件[%v]下载失败：%v", image, err)
	//	return "", err
	//}
	//content, _ := ioutil.ReadFile(local_path)
	//prefix := "data:image/jpeg;base64,"
	//if strings.HasSuffix(image, ".png") {
	//	prefix = "data:image/png;base64,"
	//} else if strings.HasSuffix(image, ".webp") {
	//	prefix = "data:image/webp;base64,"
	//}
	result := clJson.M{}
	param := clJson.M{
		"single_face": false,
		"userId":      this.UserId,
	}
	if strings.HasPrefix(image, "data:") {
		// base64
		param["img"] = image
	} else {
		param["image_url"] = image
	}
	err = client.Post(param, &result)
	// 16:22:27 akool.go:88[Err] 请求结果:map[error_code:0 error_msg:SUCCESS landmarks:[[[141 110] [189 115] [164 142] [143 163] [0 0] [0 0]]] landmarks_str:[141,110:189,115:164,142:143,163] region:[[111 58 100 132]] seconds:0.021605968475341797 trx_id:4d3673a0-6300-4951-807d-5e3e03b50d16]
	if err != nil {
		clLog.Error("访问失败:%v", err)
		return "", err
	}
	clLog.Info("请求结果:%+v", result)
	if result.Uint32("error_code") != 0 {
		return "", fmt.Errorf(result.Get("error_msg"))
	}
	//if len(result.LandmarksStr) == 0 {
	//	clLog.Error("请求结果:%+v", result)
	//	return "", fmt.Errorf("没有detect数据")
	//}
	LandmarksStr := result.Get("landmarks_str")
	LandmarksStr = strings.TrimLeft(LandmarksStr, "[")
	LandmarksStr = strings.TrimRight(LandmarksStr, "]")
	return LandmarksStr, err
}

func (this *Akool) DetectVideo(src, frame_time string) (string, error) {
	token, err := this.GenToken()
	if err != nil {
		return "", err
	}
	url := "https://faceswap.akool.com/api/v2/faceswap/material/create"
	client := xhttp.New(url)
	client.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", token),
	})
	result := clJson.M{}
	err = client.Post(clJson.M{
		"frame_time": frame_time,
		"url":        src,
		"userId":     this.UserId,
	}, &result)
	if err != nil {
		clLog.Error("访问失败:%v", err)
		return "", err
	}
	clLog.Info("请求结果:%+v", result)
	if result.Uint32("error_code") != 0 {
		return "", fmt.Errorf(result.Get("error_msg"))
	}
	LandmarksStr := result.Get("landmarks_str")
	LandmarksStr = strings.TrimLeft(LandmarksStr, "[")
	LandmarksStr = strings.TrimRight(LandmarksStr, "]")
	return LandmarksStr, err
}
